package schedule

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"on-air/auth"
	"on-air/calendarutil"
	"on-air/lifxutil"
)

const FreeBusyScope = "https://www.googleapis.com/auth/calendar.freebusy"

type Schedule struct {
	Intervals []TimeBlock
}

type State string

const (
	Busy    State = "busy"
	Free    State = "free"
	Unknown State = "unknown"
)

type TimeBlock struct {
	Start time.Time
	End   time.Time
}

type Manager struct {
	sync.RWMutex
	current               Schedule
	CredsPath             string
	TokenPath             string
	CalID                 string
	Days                  int
	LifxToken             string
	LifxLightID           string
	LifxLightLabel        string
	LifxBusyColor         string
	LifxFreeColor         string
	ReloadIntervalSeconds int
}

func (m *Manager) Update(s Schedule) {
	m.Lock()
	defer m.Unlock()
	m.current = s
}

func (m *Manager) InSchedule(t time.Time) bool {
	m.RLock()
	defer m.RUnlock()
	for _, block := range m.current.Intervals {
		if t.After(block.Start) && t.Before(block.End) {
			return true
		}
	}
	return false
}

// LoadSchedule loads free/busy from GCal for now
func (m *Manager) LoadSchedule() Schedule {
	ctx := context.Background()

	client, err := auth.GetClient(ctx, m.CredsPath, m.TokenPath, FreeBusyScope)
	if err != nil {
		log.Printf("auth client: %v", err)
		return Schedule{} // Don't exit, just return empty schedule
	}
	svc, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("calendar service: %v", err)
		return Schedule{} // Don't exit, just return empty schedule
	}
	now := time.Now().UTC()
	to := now.Add(time.Duration(m.Days) * 24 * time.Hour)

	var resp *calendar.FreeBusyResponse
	var lastErr error
	maxAttempts := 3
	backoff := 1 * time.Second
	maxBackoff := 8 * time.Second
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, lastErr = calendarutil.QueryFreeBusy(
			ctx,
			svc,
			m.CalID,
			now.Format(time.RFC3339),
			to.Format(time.RFC3339),
		)
		if lastErr == nil {
			break
		}
		// Check for 500-level error
		if apiErr, ok := lastErr.(*googleapi.Error); ok && apiErr.Code >= 500 && apiErr.Code <= 599 {
			fmt.Printf("FreeBusy query attempt %d failed with 5xx error (%d): %v. Retrying in %v...\n", attempt, apiErr.Code, apiErr, backoff)
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}
		// Not a 5xx error, break and handle as usual
		break
	}
	if lastErr != nil {
		// The query failed for some reason
		log.Printf("freebusy query: %v", lastErr)
		// Return an empty schedule to just keep the system running
		return Schedule{}
	}
	for id, cal := range resp.Calendars {
		if len(cal.Busy) == 0 {
			fmt.Printf("  %s: no busy blocks 🎉\n", id)
			continue
		}
		for _, b := range cal.Busy {
			start, err := time.Parse(time.RFC3339, b.Start)
			if err != nil {
				log.Printf("parse start time: %v", err)
				continue
			}
			end, err := time.Parse(time.RFC3339, b.End)
			if err != nil {
				log.Printf("parse end time: %v", err)
				continue
			}
			return Schedule{
				Intervals: []TimeBlock{
					{Start: start, End: end},
				},
			}
		}
	}
	return Schedule{}
}

// Reloader Worker: reload schedule based on ReloadIntervalSeconds
// If ReloadIntervalSeconds is 0, it defaults to 60 seconds.
func Reloader(m *Manager) {
	interval := time.Duration(m.ReloadIntervalSeconds) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second // fallback to 60 seconds if not set
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		m.Update(m.LoadSchedule())
		fmt.Println("Schedule reloaded")
		<-ticker.C
	}
}

// Action Represents a state change event
type Action struct {
	State State // "inside" or "outside"
	Time  time.Time
}

// ActionWorker handles REST calls
func ActionWorker(ch <-chan Action, lifxToken, lifxLightID, lifxLightLabel, lifxBusyColor, lifxFreeColor string) {
	for action := range ch {
		lc := lifxutil.NewClient(lifxToken)
		light := lifxutil.Light{ID: lifxLightID, Label: lifxLightLabel}

		if action.State == Busy {
			if err := lc.SetBusy(light, lifxBusyColor); err != nil {
				fmt.Printf("Failed to set busy state: %v\n", err)
				continue
			}
			fmt.Printf("Set busy at %s\n", action.Time.Format(time.RFC3339))
		} else if action.State == Free {
			if err := lc.SetFree(light, lifxFreeColor); err != nil {
				fmt.Printf("Failed to set free state: %v\n", err)
				continue
			}
			fmt.Printf("Set free at %s\n", action.Time.Format(time.RFC3339))
		}
	}
}

// Executor detect transitions and push to worker channel
func Executor(m *Manager, ch chan<- Action) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	currentState := Unknown

	for {
		now := time.Now()
		inSchedule := m.InSchedule(now)

		newState := Free
		if inSchedule {
			newState = Busy
		}

		// Only push events when state changes
		if newState != currentState {
			ch <- Action{State: newState, Time: now}
			currentState = newState
		}

		<-ticker.C
	}
}
