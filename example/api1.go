package main

import (
	"fmt"

	"github.com/maseer/tasks"
)

func a(ping *tasks.Ping) (interface{}, error) {
	res := []string{"k1", "k2"}
	ping.ToMultiple = true
	return res, nil
}

func api1() {
	t := tasks.New()
	t.Add(a)
	a := t.Begin([]int{1, 2})
	fmt.Println(a)
}
