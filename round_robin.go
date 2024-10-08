package main

import "sync"

type RoundRobin struct {
	locker sync.Mutex

	Current int
	Pool    []string
}

func (r *RoundRobin) Get() string {
	if len(r.Pool) == 1 {
		return r.Pool[0]
	}

	r.locker.Lock()
	defer r.locker.Unlock()

	if r.Current >= len(r.Pool) {
		r.Current %= len(r.Pool)
	}

	result := r.Pool[r.Current]
	r.Current++
	return result
}
