package monitor

import (
	"fmt"
	"regexp"
	"time"
)

type rawSnapshot map[string]string

type SnapshotEvent struct {
	Event string    `json:"event"`
	Time  time.Time `json:"time"`
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

	// Unused
	LineInteraction string `json:"lineInteraction"`
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
		re := regexp.MustCompile(`(\w+) at ([\d\/\: ]+)`)

		match := re.FindStringSubmatch(v)
		s.LastTestResult.Event = match[1]
		s.LastTestResult.Time, _ = time.Parse("2006/01/02 15:04:05", match[2])
	}

	if v, found := raw.Get("Last Power Event"); found {
		re := regexp.MustCompile(`(\w+) at ([\d\/\: ]+)`)

		match := re.FindStringSubmatch(v)
		s.LastPowerEvent.Event = match[1]
		s.LastPowerEvent.Time, _ = time.Parse("2006/01/02 15:04:05", match[2])
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

	return s
}
