package statistics

import (
	"go.uber.org/atomic"
)

type Statistics struct {
	m map[string]*atomic.Float64
}

func New(keys ...string) Statistics {
	s := Statistics{
		m: make(map[string]*atomic.Float64),
	}

	for _, key := range keys {
		s.m[key] = atomic.NewFloat64(0)
	}

	return s
}

func (s Statistics) Incr(key string, i float64) {
	s.m[key].Add(i)
}

func (s Statistics) Export() map[string]float64 {
	m := make(map[string]float64)

	for k, v := range s.m {
		m[k] = v.Load()
	}

	return m
}
