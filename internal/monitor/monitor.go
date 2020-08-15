package monitor

import (
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var dataRegexp *regexp.Regexp

type Monitor struct {
	interval  time.Duration
	engaged   bool
	ticker    *time.Ticker
	snapshots []Snapshot
}

func init() {
	dataRegexp = regexp.MustCompile(`\s{2,}?([\w ]+)\.{2,} (.+)`)
}

func New(interval time.Duration) *Monitor {
	return &Monitor{
		interval: interval,
	}
}

func (m *Monitor) Start() {
	m.ticker = time.NewTicker(m.interval)
	m.engaged = true

	// Immediately grab a snapshot
	m.exec()

	go func() {
		for range m.ticker.C {
			if m.engaged {
				m.exec()
			}
		}
	}()
}

func (m *Monitor) Stop() {
	m.ticker.Stop()
	m.engaged = false
	m.ticker = nil
}

func (m *Monitor) Active() bool {
	return m.engaged
}

func (m *Monitor) exec() {
	s := m.getUPSSnapshot()

	m.snapshots = append([]Snapshot{s}, m.snapshots...)

	// Keep last week's worth of snapshots
	if maxSize, sliceSize := int64(7*24*time.Hour/m.interval), int64(len(m.snapshots)); sliceSize > maxSize {
		m.snapshots = m.snapshots[0:maxSize]
	}
}

func (m *Monitor) getUPSSnapshot() Snapshot {
	cmd := exec.Command("pwrstat", "-status")

	out, _ := cmd.Output()

	matches := dataRegexp.FindAllStringSubmatch(string(out), -1)

	raw := rawSnapshot{}
	for _, match := range matches {
		raw[strings.TrimSpace(match[1])] = match[2]
	}

	snapshot := NewFromRawSnapshot(raw)

	return snapshot
}

func (m *Monitor) GetRecentSnapshots() SnapshotSlice {
	return SnapshotSlice(m.snapshots)
}

func (m *Monitor) GetInterval() time.Duration {
	return m.interval
}
