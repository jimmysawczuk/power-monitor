package monitor

import (
	"encoding/json"
	"fmt"
	// "math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type rawSnapshot map[string]string

func (r rawSnapshot) Get(key string) (string, bool) {
	if v, exists := r[key]; exists {
		return v, true
	} else {
		return "", false
	}
}

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
}

type SnapshotAverage struct {
	// The average battery capacity remaining
	BatteryRemaining float64

	// The average load, in watts
	Load float64

	// The average battery capacity, in watts
	BatteryCapacity float64

	// The average output voltage, in volts
	OutputVoltage float64

	// The average utility voltage, in volts
	UtilityVoltage float64

	// The remaining runtime, in minutes
	RemainingRuntime float64

	// The most recent timestamp associated with the snapshots in this average
	Timestamp time.Time

	// The interval
	Interval time.Duration

	// How many data points this snapshot represents if it's a rolling average.
	AverageOf int

	Timestamps []time.Time
}

type float64prec struct {
	val       float64
	precision float64
	fmtstr    string
}

func (f float64prec) MarshalJSON() ([]byte, error) {
	v := f.val //math.Floor(f.val*(1/f.precision)+0.5) / f.precision
	return []byte(fmt.Sprintf(f.fmtstr, v)), nil
}

func (s SnapshotAverage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		BatteryRemaining float64prec `json:"batteryRemaining"`
		Load             float64prec `json:"load"`
		BatteryCapacity  float64prec `json:"batteryCapacity"`
		OutputVoltage    float64prec `json:"outputVoltage"`
		UtilityVoltage   float64prec `json:"utilityVoltage"`
		RemainingRuntime float64prec `json:"remainingRuntime"`

		Timestamp  time.Time     `json:"timestamp"`
		Interval   time.Duration `json:"interval"`
		AverageOf  int           `json:"averageOf"`
		Timestamps []time.Time   `json:"timestamps"`
	}{
		BatteryRemaining: float64prec{val: s.BatteryRemaining, precision: 1e-4, fmtstr: "%0.4f"},
		Load:             float64prec{val: s.Load, precision: 1e-4, fmtstr: "%0.4f"},
		BatteryCapacity:  float64prec{val: s.BatteryCapacity, precision: 1e-4, fmtstr: "%0.4f"},
		OutputVoltage:    float64prec{val: s.OutputVoltage, precision: 1e-4, fmtstr: "%0.4f"},
		UtilityVoltage:   float64prec{val: s.UtilityVoltage, precision: 1e-4, fmtstr: "%0.4f"},
		RemainingRuntime: float64prec{val: s.RemainingRuntime, precision: 1e-4, fmtstr: "%0.4f"},

		Timestamp:  s.Timestamp,
		Interval:   s.Interval,
		AverageOf:  s.AverageOf,
		Timestamps: s.Timestamps,
	})
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
