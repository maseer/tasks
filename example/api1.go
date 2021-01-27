package main

import (
	"errors"
	"fmt"

	"github.com/maseer/tasks"
)

func a0(data interface{}, ping *tasks.Ping) (interface{}, error) {
	i := data.(int)
	res := []string{}
	for _, v := range []int{1, 2} {
		res = append(res, fmt.Sprintf("a%d_a%d", i, v))
	}
	ping.ToMultiple = true
	return res, nil
}
func a1(data interface{}, ping *tasks.Ping) (interface{}, error) {
	s := data.(string)
	return s, nil
}

func a2(data interface{}, ping *tasks.Ping) (interface{}, error) {
	s := data.(string)
	if s == `a7_a2` {
		return nil, errors.New("a2 random error")
	}
	r := fmt.Sprintf("%s_c", s)

	fmt.Printf("__________%s", r)

	return nil, nil
}

func api1() {
	t := tasks.New()
	t.Add(a0, a1, a2)
	t.Begin([]int{5, 7, 8})
}
