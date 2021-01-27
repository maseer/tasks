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
	ps := handlePings(ping, ping.DataStart)
	t.next(ps)
}
