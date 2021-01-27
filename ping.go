package tasks

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

type Ping struct {
	Result     *Result
	DataStart  interface{}
	DataEnd    interface{}
	ToMultiple bool
	Level      int
	Error      error `json:"-"`
	HasError   bool
}

func newPing(data interface{}, result map[string]interface{}) *Ping {
	return &Ping{
		DataStart: data,
		Result:    &Result{},
		Level:     -1,
	}
}

func cloneResult(src *Result) *Result {
	dstMap := make(map[string]interface{})
	for k, v := range src.data {
		dstMap[k] = v
	}
	r := &Result{
		data: dstMap,
	}
	return r
}

func (p *Ping) Index() string {
	bs, _ := json.Marshal(p.DataStart)
	md := md5.New()
	s := md.Sum(bs)
	r := fmt.Sprintf("%d_%x", p.Level, s)
	return r
}

func (p *Ping) clone(data interface{}) *Ping {
	clone := &Ping{
		Result:    cloneResult(p.Result),
		DataStart: data,
		Level:     p.Level + 1,
	}
	return clone
}
