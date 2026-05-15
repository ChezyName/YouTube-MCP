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

// Start returns the calculated start date string
func (r *Range) Start() string {
	start, _ := r.Resolve()
	return start
}

// End returns the calculated end date string
func (r *Range) End() string {
	_, end := r.Resolve()
	return end
}

// create the data from the given string
func (r *Range) Resolve() (string, string) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// Default values
	var endDate = today
	var startDate = now.AddDate(0, 0, -30).Format("2006-01-02")

	if r == nil {
		return time.Now().AddDate(0, 0, -30).Format("2006-01-02"), time.Now().Format("2006-01-02")
	}

	input := strings.ToLower(strings.TrimSpace(string(*r)))

	//if empty should be defautls
	if input == "" {
		return startDate, endDate
	}

	switch {
	case input == "lifetime":
		return "2005-02-14", today

	case strings.Contains(input, "/"):
		parts := strings.Split(input, "/")
		if len(parts) == 2 {
			return parts[0], parts[1]
		}

	default:
		var days int
		if _, err := fmt.Sscanf(input, "%d", &days); err == nil {
			return now.AddDate(0, 0, -days).Format("2006-01-02"), today
		}
	}

	return startDate, endDate
}
