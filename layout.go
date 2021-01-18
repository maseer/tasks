package tasks

import (
	"context"
	"reflect"
)

type Handler func(ping *Ping) (interface{}, error)
type Layout struct {
	ID         uintptr
	handleFunc Handler
	pre        *Layout
	next       *Layout
	ctx        context.Context
	level      int
	task       *Task
	limit      chan int
}

func newLayout(id uintptr, t *Task, f Handler) *Layout {
	lt := &Layout{
		handleFunc: f,
		ID:         id,
		ctx:        t.ctx,
		level:      t.index,
		task:       t,
		limit:      make(chan int, t.LimitThreadNumb),
	}
	return lt
}

// func (lt *Layout) faild(ping *Ping, err error) {
// 	lt.Fin(ping, nil, err)
// }

func (lt *Layout) Fin(ping *Ping, res interface{}, err error) {
	ping.Result.Res = res
	if err != nil {
		ping.Result.Err = err
	}
	lt.task.resultChanl <- ping.Result
}

func (lt *Layout) initWatch(ping *Ping) {
	if ping.ToMultiple {
		dataF := reflect.ValueOf(ping.Data)
		lt.task.watcher.Add(0, dataF.Len())
	} else {
		lt.task.watcher.Add(0, 1)
	}
}

func (lt *Layout) isLast() bool {
	t := lt.next == nil
	return t
}

func (lt *Layout) start(p *Ping) {
	lt.firstRun(p)
}
