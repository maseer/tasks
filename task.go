package tasks

import (
	"context"
	"reflect"
)

const defaultThreadNumb = 128

type Task struct {
	ctx         context.Context
	cancle      func()
	doms        []*Layout
	index       int
	watcher     *Watcher
	resultChanl chan *Result
	preAdd      uintptr
	Limit       int
}

func New() *Task {
	ctx, cancle := context.WithCancel(context.Background())
	r := make(chan *Result)
	return &Task{
		cancle:      cancle,
		ctx:         ctx,
		watcher:     NewWatcher(r),
		resultChanl: r,
		Limit:       defaultThreadNumb,
	}
}

func (t *Task) bindPre(dom *Layout, pre uintptr) {
	for _, v := range t.doms {
		if v.ID == pre {
			dom.pre = v
			v.next = dom
		}
	}
}

func (t *Task) Add(f ...Handler) {
	for _, v := range f {
		t.addOne(v)
	}
}

func (t *Task) addOne(f Handler) {
	pt := reflect.ValueOf(f).Pointer()
	lt := NewLayout(pt, t, f)
	t.doms = append(t.doms, lt)
	t.bindPre(lt, t.preAdd)
	t.preAdd = pt
	t.index++
}

func (t *Task) startInit(data []int) {
	p := NewPing(data, make(map[string]interface{}))
	p.ToMultiple = true
	t.doms[0].start(p)
}

func (t *Task) StartInt(data []int) {
	go t.startInit(data)
}

func (t *Task) Wait() []*Result {
	res := []*Result{}
	for r1 := range t.resultChanl {
		i := r1.index
		res = append(res, r1)
		t.watcher.Done(i - 1)
		t.watcher.check()
	}
	return res
}
