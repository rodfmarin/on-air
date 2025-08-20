package schedule

import (
	"testing"
	"time"
)

func TestManagerInSchedule(t *testing.T) {
	m := &Manager{}
	start := time.Date(2025, 8, 20, 10, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	m.Update(Schedule{Intervals: []TimeBlock{{Start: start, End: end}}})

	inside := start.Add(time.Hour)
	outside := start.Add(-time.Hour)

	if !m.InSchedule(inside) {
		t.Errorf("expected %v to be inside schedule", inside)
	}
	if m.InSchedule(outside) {
		t.Errorf("expected %v to be outside schedule", outside)
	}
}

func TestManagerUpdate(t *testing.T) {
	m := &Manager{}
	start := time.Now()
	end := start.Add(time.Hour)
	m.Update(Schedule{Intervals: []TimeBlock{{Start: start, End: end}}})
	if len(m.current.Intervals) != 1 {
		t.Errorf("expected 1 interval, got %d", len(m.current.Intervals))
	}
}

func TestExecutorStateTransitions(t *testing.T) {
	m := &Manager{}
	ch := make(chan Action, 10)
	start := time.Now().Add(-time.Minute)
	end := time.Now().Add(time.Minute)
	m.Update(Schedule{Intervals: []TimeBlock{{Start: start, End: end}}})

	// Run Executor in a goroutine
	go func() {
		Executor(m, ch)
	}()

	// Wait for a state transition
	action := <-ch
	if action.State != Busy && action.State != Free {
		t.Errorf("unexpected state: %v", action.State)
	}
}
