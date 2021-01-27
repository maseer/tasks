package main

import (
	"fmt"

	"github.com/maseer/tasks"
)

func a(data interface{}, ping *tasks.Ping) (interface{}, error) {
	i := data.(int)
	res := []string{}
	for _, v := range []int{1, 2} {
		res = append(res, fmt.Sprintf("a%d_a%d", i, v))
	}
	ping.ToMultiple = true
	return res, nil
}
func b(data interface{}, ping *tasks.Ping) (interface{}, error) {
	s := data.(string)
	return s, nil
}
func api1() {
	t := tasks.New()
	t.Add(a, b)
	a := t.Begin([]int{5, 7, 8})
	fmt.Println(a)
}
