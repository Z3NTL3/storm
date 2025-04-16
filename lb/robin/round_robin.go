package robin

import (
	"iter"
	"sync/atomic"
)

type RoundRobin[T any] struct {
	set  []T
	next uint64
}

func New[T any](set ...T) *RoundRobin[T] {
	return &RoundRobin[T]{
		set:  set,
		next: 0,
	}
}

func (r *RoundRobin[T]) Next() *T {
	if len(r.set) == 0 {
		return nil
	}
	n := atomic.AddUint64(&r.next, 1)
	// n-1 actually gives us the current index
	idx := int(n-1) % len(r.set)

	return &r.set[idx]
}

func (r *RoundRobin[T]) Iter() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for range len(r.set) {
			if !yield(r.Next()) {
				return
			}
		}
	}
}
