package tasks

type Handler func(data interface{}, ping *Ping) (interface{}, error)

func (t *Task) upateWatcher(ps []*Ping) {
	w := t.watcher
	if ps == nil {
		w.Done(len(t.doms))
		return
	}
	if len(ps) > 0 {
		w.Add(ps[0].Index, len(ps))
	}
}

func (t *Task) next(ps []*Ping) {
	t.upateWatcher(ps)
	for _, v := range ps {
		go t.runPing(v)
	}
}

func (t *Task) runPing(ping *Ping) {
	lt := t.doms[ping.Index]
	resData, err := lt.runHandle(ping)
	t.Fin(ping, resData, err)
	if lt.next == nil || err != nil {
		return
	}
	ps := handlePings(ping, resData)
	t.next(ps)
}
