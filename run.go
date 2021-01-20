package tasks

import (
	"errors"
	"reflect"
)

func copyPings(ping *Ping, data interface{}) ([]*Ping, error) {
	ps := []*Ping{}
	dataF := reflect.ValueOf(data)
	if dataF.Kind() != reflect.Slice {
		return nil, errors.New("handleMultiple error, data type must be slice")
	}
	for i := 0; i < dataF.Len(); i++ {
		d := dataF.Index(i)
		n := ping.clone(d.Interface())
		ps = append(ps, n)
	}
	return ps, nil
}

func (lt *Layout) handlePings(ping *Ping, data interface{}) ([]*Ping, error) {
	// if lt.next == nil {
	// 	return nil, nil
	// }
	if !ping.ToMultiple {
		return []*Ping{ping.clone(data)}, nil
	}
	ps, err := copyPings(ping, data)
	return ps, err
}

func (lt *Layout) firstRun(ping *Ping) error {
	ps, err := lt.handlePings(ping, ping.Data)
	if err != nil {
		return err
	}
	next(ps, lt)
	return nil
}
func next(ps []*Ping, next *Layout) {
	next.upateWatcher(len(ps))
	for _, v := range ps {
		go run(v, next)
	}
}

func (lt *Layout) upateWatcher(l int) {
	w := lt.task.watcher
	if l == -1 {
		w.Done(lt.level)
		return
	}
	w.Add(lt.level, l)
}

func (lt *Layout) runHandle(ping *Ping) (interface{}, error) {
	lt.limit <- 0
	a, err := lt.handleFunc(ping)
	<-lt.limit
	return a, err
}

func run(ping *Ping, lt *Layout) {
	resData, err := lt.runHandle(ping)
	lt.Fin(ping, resData, err)
	if lt.next == nil || err != nil {
		return
	}
	ps, err := lt.handlePings(ping, resData)
	if err != nil {
		panic(err)
	}
	if lt.next != nil {
		next(ps, lt.next)
	}
}
