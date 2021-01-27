package tasks

import "sync"

type Result struct {
	data     map[string]interface{}
	Err      error
	initOnce sync.Once
	index    int
}

func (r *Result) Init() {
	r.data = make(map[string]interface{})
}

func (r *Result) Set(key string, data interface{}) {
	if r.data == nil {
		r.initOnce.Do(r.Init)
	}
	r.data[key] = data
}

func (r *Result) Get(key string) interface{} {
	if r.data == nil {
		r.initOnce.Do(r.Init)
	}
	v, ok := r.data[key]
	if !ok {
		return nil
	}
	return v
}
