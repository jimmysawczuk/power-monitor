package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var dataRegexp *regexp.Regexp

type Monitor struct {
	interval time.Duration
	engaged  bool
	ticker   *time.Ticker
}

func init() {
	dataRegexp = regexp.MustCompile(`\s{2,}?([\w ]+)\.{2,} (.+)`)
}

func New(interval time.Duration) Monitor {
	m := Monitor{
		interval: interval,
	}

	return m
}

func (m *Monitor) Start() {
	m.ticker = time.NewTicker(m.interval)
	m.engaged = true

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

func (m Monitor) Exec() {
	log.Println("Logging")
	cmd := exec.Command("pwrstat", "-status")

	out, _ := cmd.Output()

	matches := dataRegexp.FindAllStringSubmatch(string(out), -1)

	raw := rawSnapshot{}
	for _, match := range matches {
		raw[strings.TrimSpace(match[1])] = match[2]
	}

	snapshot := NewFromRawSnapshot(raw)

	str, _ := json.MarshalIndent(snapshot, "", "   ")
	fmt.Println(string(str))
}
