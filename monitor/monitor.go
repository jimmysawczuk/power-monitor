package monitor

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var dataRegexp *regexp.Regexp

type Monitor struct {
	Interval         time.Duration
	engaged          bool
	ticker           *time.Ticker
	recent_snapshots []Snapshot
}

func init() {
	dataRegexp = regexp.MustCompile(`\s{2,}?([\w ]+)\.{2,} (.+)`)
	_ = log.Printf
}

func New(interval time.Duration) Monitor {
	m := Monitor{
		Interval: interval,
	}

	return m
}

func (m *Monitor) Start() {
	m.ticker = time.NewTicker(m.Interval)
	m.engaged = true

	// Immediately grab a snapshot
	m.Exec()

	go func() {
		for range m.ticker.C {
			if m.engaged {
				m.Exec()
			}
		}
	}()
}

func (m *Monitor) Stop() {
	m.ticker.Stop()
	m.engaged = false
	m.ticker = nil
}

func (m Monitor) Active() bool {
	return m.engaged
}

func (m *Monitor) Exec() {
	// log.Println("Logging")
	s := m.getUPSSnapshot()

	m.recent_snapshots = append([]Snapshot{s}, m.recent_snapshots...)
	if max_size, slice_size := int64(48*time.Hour/m.Interval), int64(len(m.recent_snapshots)); slice_size > max_size {
		m.recent_snapshots = m.recent_snapshots[0:max_size]
	}
}

func (m Monitor) getUPSSnapshot() Snapshot {
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

func (m Monitor) GetRecentSnapshots() SnapshotSlice {
	return SnapshotSlice(m.recent_snapshots)
}
