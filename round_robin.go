package main

import "sync"

type RoundRobin struct {
	sync.Mutex

	Current int
	Pool    []string
}

func (r *RoundRobin) Get() string {
	r.Lock()
	defer r.Unlock()

	if r.Current >= len(r.Pool) {
		r.Current %= len(r.Pool)
	}

	result := r.Pool[r.Current]
	r.Current++
	return result
}
