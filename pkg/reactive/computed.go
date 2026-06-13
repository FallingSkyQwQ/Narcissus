package reactive

import (
	"sync"
)

// Computed represents a derived reactive value.
// This is a simplified implementation where the caller must manually call Recompute().
type Computed[T any] struct {
	mu      sync.RWMutex
	value   T
	compute func() T
}

// NewComputed creates a computed signal.
// The compute function is called immediately to get the initial value.
// Note: This simplified implementation requires manual dependency tracking.
// The caller should call Recompute() when source signals change.
func NewComputed[T any](compute func() T) *Computed[T] {
	return &Computed[T]{
		value:   compute(),
		compute: compute,
	}
}

// Get returns the current computed value
func (c *Computed[T]) Get() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// SetCompute updates the compute function and recalculates the value
func (c *Computed[T]) SetCompute(compute func() T) {
	c.mu.Lock()
	c.compute = compute
	c.value = compute()
	c.mu.Unlock()
}

// Recompute recalculates the value using the current compute function
func (c *Computed[T]) Recompute() {
	c.mu.Lock()
	c.value = c.compute()
	c.mu.Unlock()
}

// Subscribe registers an observer that will be called with the current value.
// Note: In this simplified implementation, the observer is only called immediately.
// The caller is responsible for monitoring changes and calling Recompute().
func (c *Computed[T]) Subscribe(observer func(T)) func() {
	c.mu.RLock()
	value := c.value
	c.mu.RUnlock()

	// Call observer immediately with current value
	observer(value)

	// Return a no-op unsubscribe function for this simplified implementation
	return func() {}
}

// Dispose cleans up the computed signal (no-op in this simplified implementation)
func (c *Computed[T]) Dispose() {
	// No-op in simplified implementation
}
