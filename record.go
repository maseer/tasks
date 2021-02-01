package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const recordFile = `save`

type Record struct {
	E bool                   //has error
	S interface{}            //data start
	M map[string]interface{} //result
	R interface{}            //data end
}

var sepRecord = []byte(`|`)

func (r *Record) Encode() []byte {
	if r.E {
		return bytes.Join([][]byte{mustToS(r.S), mustToS(r.M)}, sepRecord)
	}
	return bytes.Join([][]byte{mustRtoS(r.R)}, sepRecord)
}

func Decode(s []byte) (*Record, error) {
	rd := &Record{}
	span := bytes.Split(s, sepRecord)
	if len(span) == 1 {
		rd.R = span[0]
	} else {
		rd.E = true
		var rmp map[string]interface{}
		err := json.Unmarshal(span[1], &rmp)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		rd.M = rmp
	}
	return rd, nil
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

	rm := make(map[string]*Record)

	bs, err := ioutil.ReadFile(recordFile)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(bs, []byte("\n"))
	for _, v := range lines {
		two := bytes.Split(v, []byte("::"))
		if len(two) != 2 {
			continue
		}
		d1, err := Decode(two[1])
		if err != nil {
			return nil, err
		}
		rm[fmt.Sprintf("%s", two[0])] = d1
	}
	return rm, nil
}

var saveLock sync.Mutex

func (mr *RecordMap) save() error {
	saveLock.Lock()
	defer saveLock.Unlock()
	fi, err := os.OpenFile(recordFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer fi.Close()

	for k, v := range *mr {
		fi.Write([]byte(fmt.Sprintf("%s::", k)))
		fi.Write(v.Encode())
		fi.Write([]byte("\n"))

	}
	return nil
}

func (t *Task) record(p *Ping) {
	if !t.cfg.UseRecord {
		return
	}
	dataRecordLock.Lock()
	defer dataRecordLock.Unlock()

	dataRecord[p.Index()] = p.toRecord(p.DataEnd)
	if err := dataRecord.save(); err != nil {
		fmt.Println(err)
	}
}

func mustRtoS(v interface{}) []byte {
	if v == nil {
		return []byte{}
	}
	return []byte(fmt.Sprintf("%s", v))
}

func mustToS(v interface{}) []byte {
	bs, err := json.Marshal(v)
	if err != nil {
		return []byte{}
	}
	return bs
}
