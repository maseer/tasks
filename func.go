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

func runFromRecord(ping *Ping) (*Ping, bool) {
	dataRecordLock.Lock()
	defer dataRecordLock.Unlock()
	index := ping.Index()
	f, ok := dataRecord[index]
	return f, ok
}

func (t *Task) runHandle(ping *Ping) (interface{}, error) {
	nping, ok := runFromRecord(ping)
	if t.UseRecord && ok && !nping.HasError {
		ping.ToMultiple = nping.ToMultiple
		ping.Result = nping.Result
		ping.DataEnd = nping.DataEnd
		return nping.DataEnd, nil
	}

	lt := t.doms[ping.Level]
	lt.limit <- 0
	defer func() {
		<-lt.limit
	}()
	i, err := lt.handleFunc(ping.DataStart, ping)
	ping.DataEnd = i
	return i, err
}
