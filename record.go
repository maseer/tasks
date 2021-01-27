package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const recordFile = `save.json`

type Record map[string]*Ping

var dataRecord = mustReadRecord()
var dataRecordLock sync.Mutex

func newRecord() Record {
	r1 := make(map[string]*Ping)
	r2 := Record(r1)
	return r2
}

func mustReadRecord() Record {
	r, _ := readRecord()
	return r
}

func readRecord() (Record, error) {
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
	var r Record
	if err := json.NewDecoder(fi).Decode(&r); err != nil {
		return newRecord(), err
	}
	return r, nil
}

var saveLock sync.Mutex

func (r *Record) save() error {
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
	if !t.UseRecord {
		return
	}

	dataRecordLock.Lock()
	defer dataRecordLock.Unlock()

	dataRecord[p.Index()] = p
	if err := dataRecord.save(); err != nil {
		fmt.Println(err)
	}
}
