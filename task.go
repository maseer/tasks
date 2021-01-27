package tasks

import (
	"context"
	"reflect"
	"time"
)

const defaultThreadNumb = 128

type Task struct {
	ctx             context.Context
	cancle          func()
	doms            []*Layout
	index           int
	watcher         *watcher
	resultChanl     chan *Ping
	preAdd          uintptr
	LimitThreadNumb int
	UseRecord       bool
}

func New() *Task {
	ctx, cancle := context.WithCancel(context.Background())
	r := make(chan *Ping)
	return &Task{
		cancle:          cancle,
		ctx:             ctx,
		watcher:         newWatcher(r),
		resultChanl:     r,
		LimitThreadNumb: defaultThreadNumb,
		UseRecord:       true,
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
		t.start(p)
	} else {
		t.resultChanl <- p
	}
}

func (t *Task) wait() []*Result {
	res := []*Result{}
	for p := range t.resultChanl {
		res = append(res, p.Result)
		t.watcher.Done(p.Level)
		t.End(p)
		if p.Level >= len(t.doms)-1 {
			<-time.After(time.Microsecond * 10) //TODO remove time
			t.watcher.check()
		}
	}
	return res
}

func (t *Task) Begin(data interface{}) []*Result {
	go t.startInit(data)
	return t.wait()
}

func (t *Task) BeginInt(data []int) []*Result {
	go t.startInit(data)
	return t.wait()
}
func (t *Task) BeginString(data []string) []*Result {
	go t.startInit(data)
	return t.wait()
}
