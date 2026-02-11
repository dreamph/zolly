package main

import "sync"

type RoundRobin struct {
	mu      sync.Mutex
	current int
	Pool    []string
}

func (r *RoundRobin) Get() string {
	if len(r.Pool) == 1 {
		return r.Pool[0]
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.current >= len(r.Pool) {
		r.current %= len(r.Pool)
	}

	result := r.Pool[r.current]
	r.current++
	return result
}
