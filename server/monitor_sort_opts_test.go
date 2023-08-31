package server

import (
	"sort"
	"testing"
	"time"
)

func TestSortByIdleTime(t *testing.T) {
	now := time.Now().UTC()

	cases := map[string]ConnInfos{
		"zero values": {{}, {}, {}, {}},
		"equal last activity times": {
			{Start: now.Add(-50 * time.Minute), LastActivity: now.Add(-time.Minute)},
			{Start: now.Add(-30 * time.Minute), LastActivity: now.Add(-time.Minute)},
			{Start: now.Add(-10 * time.Second), LastActivity: now.Add(-time.Minute)},
			{Start: now.Add(-2 * time.Hour), LastActivity: now.Add(-time.Minute)},
		},
		"last activity in the future": {
			{Start: now.Add(-50 * time.Minute), LastActivity: now.Add(10 * time.Minute)}, // +10m
			{Start: now.Add(-30 * time.Minute), LastActivity: now.Add(5 * time.Minute)},  // +5m
			{Start: now.Add(-24 * time.Hour), LastActivity: now.Add(2 * time.Second)},    // +2s
			{Start: now.Add(-10 * time.Second), LastActivity: now.Add(15 * time.Minute)}, // +15m
			{Start: now.Add(-2 * time.Hour), LastActivity: now.Add(time.Minute)},         // +1m
		},
		"unsorted": {
			{Start: now.Add(-50 * time.Minute), LastActivity: now.Add(-10 * time.Minute)}, // 10m ago
			{Start: now.Add(-30 * time.Minute), LastActivity: now.Add(-5 * time.Minute)},  // 5m ago
			{Start: now.Add(-24 * time.Hour), LastActivity: now.Add(-2 * time.Second)},    // 2s ago
			{Start: now.Add(-10 * time.Second), LastActivity: now.Add(-15 * time.Minute)}, // 15m ago
			{Start: now.Add(-2 * time.Hour), LastActivity: now.Add(-time.Minute)},         // 1m ago
		},
		"unsorted with zero value start time": {
			{LastActivity: now.Add(-10 * time.Minute)}, // 10m ago
			{LastActivity: now.Add(-5 * time.Minute)},  // 5m ago
			{LastActivity: now.Add(-2 * time.Second)},  // 2s ago
			{LastActivity: now.Add(-15 * time.Minute)}, // 15m ago
			{LastActivity: now.Add(-time.Minute)},      // 1m ago
		},
		"sorted": {
			{Start: now.Add(-24 * time.Hour), LastActivity: now.Add(-2 * time.Second)},    // 2s ago
			{Start: now.Add(-2 * time.Hour), LastActivity: now.Add(-time.Minute)},         // 1m ago
			{Start: now.Add(-30 * time.Minute), LastActivity: now.Add(-5 * time.Minute)},  // 5m ago
			{Start: now.Add(-50 * time.Minute), LastActivity: now.Add(-10 * time.Minute)}, // 10m ago
			{Start: now.Add(-10 * time.Second), LastActivity: now.Add(-15 * time.Minute)}, // 15m ago
		},
		"sorted with zero value start time": {
			{LastActivity: now.Add(-2 * time.Second)},  // 2s ago
			{LastActivity: now.Add(-time.Minute)},      // 1m ago
			{LastActivity: now.Add(-5 * time.Minute)},  // 5m ago
			{LastActivity: now.Add(-10 * time.Minute)}, // 10m ago
			{LastActivity: now.Add(-15 * time.Minute)}, // 15m ago
		},
	}

	for name, conns := range cases {
		t.Run(name, func(t *testing.T) {
			sort.Sort(byIdle{conns, now})

			idleDurations := getIdleDurations(conns, now)

			if !sortedDurationsAsc(idleDurations) {
				t.Errorf("want durations sorted in ascending order, got %v", idleDurations)
			}
		})
	}
}

// getIdleDurations returns a slice of idle durations from a connection info list up until now time.
func getIdleDurations(conns ConnInfos, now time.Time) []time.Duration {
	durations := make([]time.Duration, 0, len(conns))

	for _, conn := range conns {
		durations = append(durations, now.Sub(conn.LastActivity))
	}

	return durations
}

// sortedDurationsAsc checks if a time.Duration slice is sorted in ascending order.
func sortedDurationsAsc(durations []time.Duration) bool {
	return sort.SliceIsSorted(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
}
