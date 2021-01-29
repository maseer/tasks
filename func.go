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

func tryRecord(ping *Ping) (*Record, bool) {
	dataRecordLock.Lock()
	defer dataRecordLock.Unlock()
	index := ping.Index()
	f, ok := dataRecord[index]
	return f, ok
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
	if t.UseRecord && t.isLast(ping) {
		r, ok := tryRecord(ping)
		if ok && !r.E {
			return r.R, nil
		}
	}
	lt := t.doms[ping.Level]
	lt.limit <- 0
	defer func() {
		<-lt.limit
	}()
	i, err := lt.handleFunc(ping.DataStart, ping)
	return i, err
}
