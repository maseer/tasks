package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const recordFile = `save.json`

type Record struct {
	D interface{}            //data
	E bool                   //has error
	M map[string]interface{} //*result
	R interface{}            //data resultd
}

type RecordMap map[string]*Record

var dataRecord = mustReadRecord()
var dataRecordLock sync.Mutex

func newRecord() RecordMap {
	r1 := make(map[string]*Record)
	r2 := RecordMap(r1)
	return r2
}

func mustReadRecord() RecordMap {
	r, _ := readRecord()
	return r
}

func readRecord() (RecordMap, error) {
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
	var r RecordMap
	if err := json.NewDecoder(fi).Decode(&r); err != nil {
		return newRecord(), err
	}
	return r, nil
}

var saveLock sync.Mutex

func (r *RecordMap) save() error {
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

func (t *Task) record(p *Ping) {
	if !t.UseRecord {
		return
	}
	dataRecordLock.Lock()
	defer dataRecordLock.Unlock()

	dataRecord[p.Index()] = p.toRecord(p.DataEnd)
	if err := dataRecord.save(); err != nil {
		fmt.Println(err)
	}
}
