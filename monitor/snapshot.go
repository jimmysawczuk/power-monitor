package monitor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type rawSnapshot map[string]string

type SnapshotEvent struct {
	Event    string        `json:"event"`
	Time     time.Time     `json:"time"`
	Duration time.Duration `json:"duration"`
}

type Snapshot struct {
	// The current battery capacity remaining
	BatteryRemaining float64 `json:"batteryRemaining"`

	// The current load, in watts
	Load int64 `json:"load"`

	// The battery capacity, in watts
	BatteryCapacity int64 `json:"batteryCapacity"`

	// The current output voltage, in volts
	OutputVoltage int64 `json:"outputVoltage"`

	// The current utility voltage, in volts
	UtilityVoltage int64 `json:"utilityVoltage"`

	// The remaining runtime, in minutes
	RemainingRuntime int64 `json:"remainingRuntime"`

	// Status: "Normal", others are options
	Status string `json:"status"`

	// Last test result
	LastTestResult SnapshotEvent `json:"lastTestResult"`

	// Last power event
	LastPowerEvent SnapshotEvent `json:"lastPowerEvent"`

	// The current firmware version, model name
	FirmwareVersion string `json:"firmwareVersion"`
	ModelName       string `json:"modelName"`

	// The timestamp on this snapshot
	Timestamp time.Time `json:"timestamp"`

	// Unused
	LineInteraction string `json:"lineInteraction"`

	// How many data points this snapshot represents if it's a rolling average.
	AverageOf int `json:"averageOf,omitempty"`
}

func (r rawSnapshot) Get(key string) (string, bool) {
	if v, exists := r[key]; exists {
		return v, true
	} else {
		return "", false
	}
}

func NewFromRawSnapshot(raw rawSnapshot) (s Snapshot) {
	if v, found := raw.Get("Battery Capacity"); found {
		fmt.Sscanf(v, "%f %%", &s.BatteryRemaining)
		s.BatteryRemaining /= 100.0
	}

	if v, found := raw.Get("Load"); found {
		fmt.Sscanf(v, "%d Watt(%d %%)", &s.Load)
	}

	if v, found := raw.Get("Rating Power"); found {
		fmt.Sscanf(v, "%d Watt", &s.BatteryCapacity)
	}

	if v, found := raw.Get("Output Voltage"); found {
		fmt.Sscanf(v, "%d V", &s.OutputVoltage)
	}

	if v, found := raw.Get("Utility Voltage"); found {
		fmt.Sscanf(v, "%d V", &s.UtilityVoltage)
	}

	if v, found := raw.Get("Remaining Runtime"); found {
		fmt.Sscanf(v, "%d min.", &s.RemainingRuntime)
	}

	if v, found := raw.Get("State"); found {
		fmt.Sscanf(v, "%s", &s.Status)
	}

	if v, found := raw.Get("Test Result"); found {
		re := regexp.MustCompile(`([\w ]+) at ([\d\/\: ]+)`)

		if re.MatchString(v) {
			match := re.FindStringSubmatch(v)
			s.LastTestResult.Event = match[1]
			s.LastTestResult.Time, _ = time.Parse("2006/01/02 15:04:05", match[2])
		} else if strings.ToLower(v) == "in progress" {
			s.LastTestResult.Event = "In Progress"
			s.LastTestResult.Time = time.Now()
		} else if strings.ToLower(v) == "unknown" {
			s.LastTestResult.Event = "Unknown"
		}
	}

	if v, found := raw.Get("Last Power Event"); found {
		re := regexp.MustCompile(`([\w ]+) at ([\d\/\: ]+)( ?for (\d+) ([\w\.]+))?`)

		if re.MatchString(v) {
			match := re.FindStringSubmatch(v)
			s.LastPowerEvent.Event = match[1]
			s.LastPowerEvent.Time, _ = time.Parse("2006/01/02 15:04:05", strings.TrimSpace(match[2]))

			if match[4] != "" {
				dur, err := strconv.ParseInt(match[4], 10, 64)
				if err == nil {
					switch match[5] {
					case "sec.":
						s.LastPowerEvent.Duration = time.Duration(dur) * time.Second
					case "min.":
						s.LastPowerEvent.Duration = time.Duration(dur) * time.Minute
					}
				}
			}
		}
	}

	if v, found := raw.Get("Model Name"); found {
		fmt.Sscanf(v, "%s", &s.ModelName)
	}

	if v, found := raw.Get("Firmware Number"); found {
		fmt.Sscanf(v, "%s", &s.FirmwareVersion)
	}

	if v, found := raw.Get("Line Interaction"); found {
		fmt.Sscanf(v, "%s", &s.LineInteraction)
	}

	s.Timestamp = time.Now().Round(time.Second)

	return s
}
