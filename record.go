package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const recordFile = `save.json`

type record map[string]*Ping

func newRecord() record {
	r1 := make(map[string]*Ping)
	r2 := record(r1)
	return r2
}

func readRecord() (record, error) {
	saveLock.Lock()
	defer saveLock.Unlock()
	if _, err := os.Stat(recordFile); err != nil {
		return newRecord(), err
	}
	fi, err := os.OpenFile(recordFile, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return newRecord(), err
	}
	defer fi.Close()
	var r record
	if err := json.NewDecoder(fi).Decode(&r); err != nil {
		return newRecord(), err
	}
	return r, nil
}

var saveLock sync.Mutex

func (r *record) save() error {
	saveLock.Lock()
	defer saveLock.Unlock()
	fi, err := os.OpenFile(recordFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	defer fi.Close()
	if err := json.NewEncoder(fi).Encode(r); err != nil {
		return err
	}
	return nil
}

func (t *Task) End(p *Ping) {
	r, err := readRecord()
	if err != nil {
		fmt.Println(err)
	}
	r[p.Index()] = p
	if err := r.save(); err != nil {
		fmt.Println(err)
	}
}
