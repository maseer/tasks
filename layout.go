package tasks

import (
	"context"
)

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

func (t *Task) Fin(ping *Ping, res interface{}, err error) {
	if err != nil {
		ping.Result.Err = err
	}
	t.resultChanl <- ping.Result
}

func (t *Task) start(p *Ping) {
	t.firstRun(p)
}
