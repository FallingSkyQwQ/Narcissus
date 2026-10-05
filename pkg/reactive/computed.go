package reactive

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Computed represents a derived reactive value.
//
// It tracks the signals and other computeds its compute function reads, and
// re-evaluates automatically when any of them change. Evaluation is lazy: the
// value is recomputed on the next Get (or when a subscriber needs it), so a
// chain of computeds that nobody reads costs nothing.
type Computed[T any] struct {
	dep     *dependent
	compute func() T

	mu          sync.Mutex
	value       T
	initialized atomic.Bool

	subMu       sync.Mutex
	subscribers map[uint64]func(T)
	nextID      atomic.Uint64
}

// NewComputed creates a computed signal and evaluates it once for the initial
// value. Dependencies are collected automatically during that evaluation.
func NewComputed[T any](compute func() T) *Computed[T] {
	c := &Computed[T]{
		dep:         newDependent(),
		compute:     compute,
		subscribers: make(map[uint64]func(T)),
	}
	c.ensure()
	return c
}

// Get returns the current value, recomputing it first if a dependency changed.
func (c *Computed[T]) Get() T {
	if outer := currentDependent(); outer != nil {
		outer.trackSource(c.dep)
	}
	c.ensure()

	c.mu.Lock()
	value := c.value
	c.mu.Unlock()
	return value
}

// SetCompute updates the compute function and recalculates the value.
func (c *Computed[T]) SetCompute(compute func() T) {
	if compute == nil {
		return
	}
	c.mu.Lock()
	c.compute = compute
	c.mu.Unlock()
	c.Recompute()
}

// Recompute forces a recalculation using the current compute function.
func (c *Computed[T]) Recompute() {
	c.dep.mu.Lock()
	c.dep.invalid = true
	c.dep.mu.Unlock()
	c.ensure()
}

// Subscribe registers an observer and returns an unsubscribe function. The
// observer is called immediately with the current value and then whenever the
// computed value changes. While at least one subscriber exists the computed is
// evaluated eagerly; with none it stays lazy.
func (c *Computed[T]) Subscribe(observer func(T)) func() {
	if observer == nil {
		return func() {}
	}

	id := c.nextID.Add(1) - 1
	c.subMu.Lock()
	c.subscribers[id] = observer
	c.setEagerLocked(true)
	c.subMu.Unlock()

	observer(c.Get())

	return func() {
		c.subMu.Lock()
		delete(c.subscribers, id)
		if len(c.subscribers) == 0 {
			c.setEagerLocked(false)
		}
		c.subMu.Unlock()
	}
}

// Dispose releases the computed's dependencies.
func (c *Computed[T]) Dispose() {
	c.dep.dispose()

	c.subMu.Lock()
	c.subscribers = make(map[uint64]func(T))
	c.subMu.Unlock()
}

func (c *Computed[T]) ensure() {
	c.dep.mu.Lock()
	if !c.dep.invalid && c.initialized.Load() {
		c.dep.mu.Unlock()
		return
	}
	c.dep.invalid = false
	c.dep.mu.Unlock()

	c.mu.Lock()
	compute := c.compute
	c.mu.Unlock()
	if compute == nil {
		return
	}

	track(c.dep, func() {
		value := compute()

		c.mu.Lock()
		changed := !c.initialized.Load() || !reflect.DeepEqual(value, c.value)
		c.value = value
		c.mu.Unlock()
		c.initialized.Store(true)

		if changed {
			c.notify(value)
		}
	})
}

// setEagerLocked toggles whether an invalidation schedules an eager
// re-evaluation. Callers must hold subMu.
func (c *Computed[T]) setEagerLocked(eager bool) {
	c.dep.mu.Lock()
	if eager {
		c.dep.flush = c.ensure
	} else {
		c.dep.flush = nil
	}
	c.dep.mu.Unlock()
}

func (c *Computed[T]) notify(value T) {
	c.subMu.Lock()
	subscribers := make([]func(T), 0, len(c.subscribers))
	for _, observer := range c.subscribers {
		subscribers = append(subscribers, observer)
	}
	c.subMu.Unlock()

	for _, observer := range subscribers {
		observer(value)
	}
}
