package youtube

import (
	"fmt"
	"strings"
	"time"
)

type Range string

type RangeDetails struct {
	startDate string
	endDate   string
}

// Start returns the calculated start date string.
func (r *Range) Start(optionalDefault ...string) string {
	start, _ := r.Resolve(optionalDefault...)
	return start
}

// End returns the calculated end date string.
func (r *Range) End(optionalDefault ...string) string {
	_, end := r.Resolve(optionalDefault...)
	return end
}

// create the data from the given string
func (r *Range) Resolve(optionalDefault ...string) (string, string) {
	now := time.Now()
	today := now.Format("2006-01-02")
	defaultStart := now.AddDate(0, 0, -30).Format("2006-01-02")

	var fallback string
	if len(optionalDefault) > 0 {
		fallback = optionalDefault[0]
	}

	// If fallback itself is "lifetime" or "x days", resolve it against today;
	// otherwise use the 30-day default start date.
	resolvedDefaultStart, resolvedDefaultEnd := defaultStart, today
	if strings.TrimSpace(fallback) != "" {
		resolvedDefaultStart, resolvedDefaultEnd = resolveRangeString(fallback, now, today, defaultStart)
	}

	if r == nil || strings.TrimSpace(string(*r)) == "" {
		return resolvedDefaultStart, resolvedDefaultEnd
	}

	return resolveRangeString(string(*r), now, today, resolvedDefaultStart)
}

func resolveRangeString(raw string, now time.Time, today string, fallbackStart string) (string, string) {
	input := strings.ToLower(strings.TrimSpace(raw))
	if input == "" {
		return fallbackStart, today
	}

	switch {
	case input == "lifetime":
		return "2005-02-14", today

	case strings.Contains(input, "/"):
		parts := strings.Split(input, "/")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}

	default:
		var days int
		if _, err := fmt.Sscanf(input, "%d", &days); err == nil && days >= 0 {
			return now.AddDate(0, 0, -days).Format("2006-01-02"), today
		}
	}

	return fallbackStart, today
}
