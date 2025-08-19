// Package calendarutil provides utilities for interacting with the Google Calendar API.
package calendarutil

import (
	"context"

	"google.golang.org/api/calendar/v3"
)

// QueryFreeBusy queries the FreeBusy endpoint for the given calendar ID and time range.
func QueryFreeBusy(ctx context.Context, svc *calendar.Service, calID string, timeMin, timeMax string) (*calendar.FreeBusyResponse, error) {
	req := &calendar.FreeBusyRequest{
		TimeMin: timeMin,
		TimeMax: timeMax,
		Items:   []*calendar.FreeBusyRequestItem{{Id: calID}},
	}
	return svc.Freebusy.Query(req).Do()
}
