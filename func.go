package tasks

type Handler func(data interface{}, ping *Ping) (interface{}, error)
