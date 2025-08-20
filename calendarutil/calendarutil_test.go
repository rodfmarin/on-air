package calendarutil

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func getTestService(t *testing.T) *calendar.Service {
	creds := "../credentials.json"
	tokenFile := "../token.json"
	if _, err := os.Stat(creds); err != nil {
		t.Skip("credentials.json not found, skipping integration test")
	}
	if _, err := os.Stat(tokenFile); err != nil {
		t.Skip("token.json not found, skipping integration test")
	}

	// Load credentials
	credData, err := os.ReadFile(creds)
	if err != nil {
		t.Skip("unable to read credentials.json, skipping integration test")
	}
	config, err := google.ConfigFromJSON(credData, calendar.CalendarScope)
	if err != nil {
		t.Skipf("failed to parse credentials.json: %v", err)
	}

	// Load token
	tokenData, err := os.ReadFile(tokenFile)
	if err != nil {
		t.Skip("unable to read token.json, skipping integration test")
	}
	var token oauth2.Token
	if err := json.Unmarshal(tokenData, &token); err != nil {
		t.Skipf("failed to parse token.json: %v", err)
	}

	// Create HTTP client
	client := config.Client(context.Background(), &token)
	ctx := context.Background()
	svc, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		t.Skipf("failed to create calendar service: %v", err)
	}
	return svc
}

func isServiceAccount(data []byte) bool {
	return bytes.Contains(data, []byte(`"type": "service_account"`))
}

func TestQueryFreeBusy_Valid(t *testing.T) {
	svc := getTestService(t)
	calID := "primary"
	now := time.Now().UTC()
	to := now.Add(24 * time.Hour)
	resp, err := QueryFreeBusy(context.Background(), svc, calID, now.Format(time.RFC3339), to.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("QueryFreeBusy failed: %v", err)
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

func TestQueryFreeBusy_InvalidCalendarID(t *testing.T) {
	svc := getTestService(t)
	now := time.Now().UTC()
	to := now.Add(24 * time.Hour)
	resp, err := QueryFreeBusy(context.Background(), svc, "invalid_calendar_id", now.Format(time.RFC3339), to.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error for invalid calendar ID: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response for invalid calendar ID")
	}
	cal, ok := resp.Calendars["invalid_calendar_id"]
	if !ok {
		t.Errorf("expected Calendars to contain invalid_calendar_id, got: %+v", resp.Calendars)
	}
	if len(cal.Busy) != 0 {
		t.Errorf("expected no busy blocks for invalid calendar ID, got: %+v", cal.Busy)
	}
}

func TestQueryFreeBusy_InvalidTimeRange(t *testing.T) {
	svc := getTestService(t)
	calID := "primary"
	// Invalid time format
	_, err := QueryFreeBusy(context.Background(), svc, calID, "not-a-time", "not-a-time")
	if err == nil {
		t.Error("expected error for invalid time range, got nil")
	}
}
