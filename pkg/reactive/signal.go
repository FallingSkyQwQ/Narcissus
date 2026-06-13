package reactive

import (
	"sync"
	"sync/atomic"
)

// Signal represents a reactive state container
type Signal[T any] struct {
	mu        sync.RWMutex
	value     T
	observers map[uint64]func(T)
	nextID    atomic.Uint64
}

// NewSignal creates a new Signal with an initial value
func NewSignal[T any](initial T) *Signal[T] {
	return &Signal[T]{
		value:     initial,
		observers: make(map[uint64]func(T)),
	}
}

// Get returns the current value
func (s *Signal[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set updates the value and notifies observers
func (s *Signal[T]) Set(newValue T) {
	s.mu.Lock()
	s.value = newValue
	observers := make([]func(T), 0, len(s.observers))
	for _, observer := range s.observers {
		observers = append(observers, observer)
	}
	s.mu.Unlock()

	// Notify outside of lock to prevent deadlocks
	for _, observer := range observers {
		observer(newValue)
	}
}

// Subscribe registers an observer and returns an unsubscribe function
func (s *Signal[T]) Subscribe(observer func(T)) func() {
	// Generate unique ID using atomic increment (returns new value, so subtract 1 for 0-based index)
	id := s.nextID.Add(1) - 1

	s.mu.Lock()
	s.observers[id] = observer
	currentValue := s.value
	s.mu.Unlock()

	// Call observer immediately with current value
	observer(currentValue)

	return func() {
		s.mu.Lock()
		delete(s.observers, id)
		s.mu.Unlock()
	}
}
