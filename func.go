package tasks

type Handler func(i interface{}, p *Ping) (interface{}, error)

func (t *Task) upateWatcher(ps []*Ping) {
	w := t.watcher
	if ps == nil {
		w.Done(len(t.doms))
		return
	}
	if len(ps) > 0 {
		w.Add(ps[0].Level, len(ps))
	}
}

func (t *Task) next(ps []*Ping) {
	t.upateWatcher(ps)
	for _, v := range ps {
		go t.runPing(v)
	}
}

func (t *Task) runPing(ping *Ping) {
	lt := t.doms[ping.Level]
	resData, err := t.runHandle(ping)
	t.Fin(ping, resData, err)
	if lt.next == nil || err != nil {
		return
	}
	ps := handlePings(ping, resData)
	t.next(ps)
}

func (t *Task) runHandle(ping *Ping) (interface{}, error) {
	lt := t.doms[ping.Level]
	lt.limit <- 0
	a, err := lt.handleFunc(ping.Data, ping)
	<-lt.limit
	return a, err
}
