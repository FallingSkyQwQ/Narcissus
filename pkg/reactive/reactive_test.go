package reactive

import "testing"

func TestComputedChainedPropagation(t *testing.T) {
	a := NewSignal(2)
	doubled := NewComputed(func() int { return a.Get() * 2 })
	plusOne := NewComputed(func() int { return doubled.Get() + 1 })

	if plusOne.Get() != 5 {
		t.Fatalf("expected 5, got %d", plusOne.Get())
	}

	a.Set(10)
	if plusOne.Get() != 21 {
		t.Errorf("expected 21 after chained update, got %d", plusOne.Get())
	}
}

func TestEffectTracksMultipleSignals(t *testing.T) {
	a := NewSignal(1)
	b := NewSignal(1)

	runs := 0
	effect := NewEffect(func() {
		_ = a.Get()
		_ = b.Get()
		runs++
	})
	defer effect.Dispose()

	if runs != 1 {
		t.Fatalf("expected 1 run on creation, got %d", runs)
	}

	a.Set(2)
	if runs != 2 {
		t.Errorf("expected run after a change, got %d", runs)
	}

	b.Set(2)
	if runs != 3 {
		t.Errorf("expected run after b change, got %d", runs)
	}
}

// TestEffectDynamicDependencies verifies that a branch that stops reading a
// signal also stops reacting to it.
func TestEffectDynamicDependencies(t *testing.T) {
	useA := NewSignal(true)
	a := NewSignal(1)
	b := NewSignal(10)

	runs := 0
	effect := NewEffect(func() {
		if useA.Get() {
			_ = a.Get()
		} else {
			_ = b.Get()
		}
		runs++
	})
	defer effect.Dispose()

	if runs != 1 {
		t.Fatalf("expected 1 run, got %d", runs)
	}

	// b is not read while useA is true, so changing it must not re-run.
	b.Set(11)
	if runs != 1 {
		t.Errorf("effect reacted to an untracked signal, runs = %d", runs)
	}

	a.Set(2)
	if runs != 2 {
		t.Errorf("expected run after tracked signal change, got %d", runs)
	}

	// Switch the branch: b becomes tracked, a stops being tracked.
	useA.Set(false)
	if runs != 3 {
		t.Fatalf("expected run after branch switch, got %d", runs)
	}

	a.Set(3)
	if runs != 3 {
		t.Errorf("effect still reacted to dropped dependency, runs = %d", runs)
	}

	b.Set(12)
	if runs != 4 {
		t.Errorf("expected run after newly tracked signal change, got %d", runs)
	}
}

// TestEffectDiamondRunsOnce makes sure a dependency read through two paths runs
// the effect a single time per update.
func TestEffectDiamondRunsOnce(t *testing.T) {
	a := NewSignal(1)
	left := NewComputed(func() int { return a.Get() * 2 })
	right := NewComputed(func() int { return a.Get() + 1 })
	both := NewComputed(func() int { return left.Get() + right.Get() })

	runs := 0
	effect := NewEffect(func() {
		_ = both.Get()
		runs++
	})
	defer effect.Dispose()

	if runs != 1 {
		t.Fatalf("expected 1 run, got %d", runs)
	}

	a.Set(2)
	if runs != 2 {
		t.Errorf("diamond should run the effect once per update, got %d runs", runs)
	}
	if both.Get() != 7 {
		t.Errorf("expected 7, got %d", both.Get())
	}
}

func TestBatchCoalescesEffectRuns(t *testing.T) {
	s := NewSignal(0)

	runs := 0
	effect := NewEffect(func() {
		_ = s.Get()
		runs++
	})
	defer effect.Dispose()

	Batch(func() {
		s.Set(1)
		s.Set(2)
		s.Set(3)
	})

	if runs != 2 {
		t.Errorf("expected one batched run (2 total), got %d", runs)
	}
}

func TestBatchNested(t *testing.T) {
	s := NewSignal(0)

	runs := 0
	effect := NewEffect(func() {
		_ = s.Get()
		runs++
	})
	defer effect.Dispose()

	Batch(func() {
		s.Set(1)
		Batch(func() {
			s.Set(2)
		})
		s.Set(3)
	})

	if runs != 2 {
		t.Errorf("expected one run for nested batches, got %d", runs)
	}
}

func TestEffectCleanupRunsBeforeRerun(t *testing.T) {
	s := NewSignal(0)

	var order []string
	effect := NewEffect(func() {
		_ = s.Get()
		order = append(order, "run")
	})
	defer effect.Dispose()

	effect.SetCleanup(func() { order = append(order, "cleanup") })
	s.Set(1)

	if len(order) != 3 || order[0] != "run" || order[1] != "cleanup" || order[2] != "run" {
		t.Errorf("expected [run cleanup run], got %v", order)
	}
}
