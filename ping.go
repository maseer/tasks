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
	// if s, ok := p.DataStart.(string); ok {
	// 	s2 := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(s))
	// 	return fmt.Sprintf("%d_string_%s", p.Level, s2)
	// }
	// if s, ok := p.DataStart.(int); ok {
	// 	return fmt.Sprintf("%d_int_%d", p.Level, s)
	// }
	// if s, ok := p.DataStart.(float64); ok {
	// 	return fmt.Sprintf("%d_float_%f", p.Level, s)
	// }
	bs, _ := json.Marshal(p.DataStart)
	s := md5.Sum(bs)
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
func (p *Ping) toRecord(data interface{}) *Record {
	if !p.HasError {
		return &Record{
			R: data,
		}
	}
	r := &Record{
		S: p.DataStart,
		E: p.HasError,
		M: p.Result.data,
		R: data,
	}
	return r
}
