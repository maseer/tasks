package tasks

import (
	"errors"
	"reflect"
)

func mustCopyPings(ping *Ping, data interface{}) []*Ping {
	ps := []*Ping{}
	dataF := reflect.ValueOf(data)
	if dataF.Kind() != reflect.Slice {
		panic(errors.New("handleMultiple error, data type must be slice"))
	}
	for i := 0; i < dataF.Len(); i++ {
		d := dataF.Index(i)
		n := ping.clone(d.Interface())
		ps = append(ps, n)
	}
	return ps
}

func handlePings(ping *Ping, data interface{}) []*Ping {
	if !ping.ToMultiple {
		return []*Ping{ping.clone(data)}
	}
	ps := mustCopyPings(ping, data)
	return ps
}

func (t *Task) firstRun(ping *Ping) {
	ps := handlePings(ping, ping.Data)
	t.next(ps)
}

func (lt *Layout) runHandle(ping *Ping) (interface{}, error) {
	lt.limit <- 0
	a, err := lt.handleFunc(ping.Data, ping)
	<-lt.limit
	return a, err
}
