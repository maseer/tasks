package tasks

import (
	"fmt"
	"sync"
)

type watcher struct {
	conter []*stepCounter
	lock   sync.Mutex
	max    int
	ping   chan *Ping
}

type stepCounter struct {
	todo int
	fin  int
}

func newWatcher(r chan *Ping) *watcher {
	return &watcher{
		ping: r,
	}
}

func (c *watcher) init(step int) {
	if (step + 1) > c.max {
		c.conter = append(c.conter, &stepCounter{})
		c.max = step + 1
	}
}

func (c *watcher) Add(step int, i int) {
	c.lock.Lock()
	c.init(step)
	c.conter[step].todo += i
	c.lock.Unlock()
}

func (c *watcher) Done(step int) {
	c.lock.Lock()
	c.init(step)
	c.conter[step].fin += 1
	c.lock.Unlock()
}

func (c *watcher) check() {
	pnt := ``
	for _, v := range c.conter {
		pnt += fmt.Sprintf("[%d/%d] ", v.fin, v.todo)
	}
	fmt.Printf("%s\r", pnt)
	for _, v := range c.conter {
		if v.fin < v.todo {
			return
		}
	}
	fmt.Println("")
	close(c.ping)
}
