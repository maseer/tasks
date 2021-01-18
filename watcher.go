package tasks

import (
	"fmt"
	"sync"
)

type Watcher struct {
	conter []*StepCounter
	lock   sync.Mutex
	max    int
	result chan *Result
}

type StepCounter struct {
	todo int
	fin  int
}

func NewWatcher(r chan *Result) *Watcher {
	return &Watcher{
		result: r,
	}
}

func (c *Watcher) init(step int) {
	if (step + 1) > c.max {
		c.conter = append(c.conter, &StepCounter{})
		c.max = step + 1
	}
}

func (c *Watcher) Add(step int, i int) {
	c.lock.Lock()
	c.init(step)
	c.conter[step].todo += i
	c.lock.Unlock()
}

func (c *Watcher) Done(step int) {
	c.lock.Lock()
	c.init(step)
	c.conter[step].fin += 1
	c.lock.Unlock()
}

func (c *Watcher) check() {
	pnt := ``
	for _, v := range c.conter {
		pnt += fmt.Sprintf("[%d/%d] ", v.fin, v.todo)
	}
	fmt.Printf("%s\n", pnt)
	for _, v := range c.conter {
		if v.fin < v.todo {
			return
		}
	}
	close(c.result)
}
