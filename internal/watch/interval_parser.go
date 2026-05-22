package watch

import (
	"fmt"
	"strings"
	"time"
)

// ParseInterval parses a human-friendly interval string into a time.Duration.
// It accepts standard Go duration strings (e.g. "30s", "5m", "1h") as well as
// plain integer values suffixed with a unit word (e.g. "30 seconds", "5 minutes").
func ParseInterval(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("watch: interval must not be empty")
	}

	// Try standard Go duration first.
	if d, err := time.ParseDuration(s); err == nil {
		if d <= 0 {
			return 0, fmt.Errorf("watch: interval must be positive, got %q", s)
		}
		return d, nil
	}

	// Accept word-based suffixes.
	words := strings.Fields(s)
	if len(words) == 2 {
		var n int
		if _, err := fmt.Sscanf(words[0], "%d", &n); err == nil {
			unit := strings.ToLower(strings.TrimSuffix(words[1], "s")) // strip plural
			switch unit {
			case "second":
				return time.Duration(n) * time.Second, nil
			case "minute":
				return time.Duration(n) * time.Minute, nil
			case "hour":
				return time.Duration(n) * time.Hour, nil
			}
		}
	}

	return 0, fmt.Errorf("watch: unrecognised interval format %q", s)
}
