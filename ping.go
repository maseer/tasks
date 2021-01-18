package tasks

type Ping struct {
	Result     *Result
	Data       interface{}
	ToMultiple bool
}

func NewPing(data interface{}, result map[string]interface{}) *Ping {
	return &Ping{
		Data:   data,
		Result: &Result{},
	}
}

func cloneResult(src *Result) *Result {
	dstMap := make(map[string]interface{})
	for k, v := range src.data {
		dstMap[k] = v
	}
	r := &Result{
		data:  dstMap,
		Err:   src.Err,
		index: src.index + 1,
	}
	return r
}

func (p *Ping) Clone(data interface{}) *Ping {
	clone := &Ping{
		Result: cloneResult(p.Result),
		Data:   data,
	}
	return clone
}
