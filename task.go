package tasks

import (
	"context"
	"reflect"
)

const defaultThreadNumb = 128

type Task struct {
	ctx             context.Context
	cancle          func()
	doms            []*Layout
	index           int
	watcher         *watcher
	resultChanl     chan *Result
	preAdd          uintptr
	LimitThreadNumb int
}

func New() *Task {
	ctx, cancle := context.WithCancel(context.Background())
	r := make(chan *Result)
	return &Task{
		cancle:          cancle,
		ctx:             ctx,
		watcher:         newWatcher(r),
		resultChanl:     r,
		LimitThreadNumb: defaultThreadNumb,
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
	lt := newLayout(pt, t, f)
	t.doms = append(t.doms, lt)
	t.bindPre(lt, t.preAdd)
	t.preAdd = pt
	t.index++
}

func (t *Task) startInit(data interface{}) {
	p := newPing(data, make(map[string]interface{}))
	p.ToMultiple = true
	if len(t.doms) > 0 {
		t.doms[0].start(p)
	} else {
		t.resultChanl <- &Result{index: 1}
	}
}

func (t *Task) wait() []*Result {
	res := []*Result{}
	for r1 := range t.resultChanl {
		i := r1.index
		res = append(res, r1)
		t.watcher.Done(i - 1)
		t.watcher.check()
	}
	return res
}

func (t *Task) Begin(data []int) []*Result {
	go t.startInit(data)
	return t.wait()
}

func (t *Task) BeginString(data []string) []*Result {
	go t.startInit(data)
	return t.wait()
}
