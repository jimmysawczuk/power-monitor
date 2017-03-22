package monitor

import (
	"math"
	"time"
)

type SnapshotSlice []Snapshot

func (in SnapshotSlice) Filter(test func(Snapshot) bool) SnapshotSlice {
	out := make(SnapshotSlice, 0)

	for i := 0; i < len(in); i++ {
		if test(in[i]) {
			out = append(out, in[i])
		}
	}

	return out
}

func (in SnapshotSlice) Rollup(interval time.Duration) SnapshotSlice {
	if len(in) <= 1 {
		return in
	}

	out := make(SnapshotSlice, 0)
	lastIndex := 0

	for i := 1; i < len(in); i++ {
		if in[lastIndex].Timestamp.UnixNano()-in[i].Timestamp.UnixNano() >= int64(interval) {
			res, len := SnapshotSlice(in[lastIndex:i]).Average()
			if len > 0 {
				out = append(out, res)
			}
			lastIndex = i
		}
	}

	res, len := SnapshotSlice(in[lastIndex:]).Average()
	if len > 0 {
		out = append(out, res)
	}

	return out
}

func (in SnapshotSlice) Average() (Snapshot, int) {
	if len(in) == 0 {
		return Snapshot{}, 0
	}

	if len(in) == 1 {
		return in[0], 1
	}

	ret := Snapshot{}
	for i := 0; i < len(in); i++ {
		v := in[i]

		ret.BatteryRemaining += v.BatteryRemaining
		ret.Load += v.Load
		ret.BatteryCapacity += v.BatteryCapacity
		ret.OutputVoltage += v.OutputVoltage
		ret.UtilityVoltage += v.UtilityVoltage
		ret.RemainingRuntime += v.RemainingRuntime
	}

	ret.BatteryRemaining = ret.BatteryRemaining / float64(len(in))
	ret.Load = int64(math.Floor(float64(ret.Load)/float64(len(in)) + 0.5))
	ret.BatteryCapacity = int64(math.Floor(float64(ret.BatteryCapacity)/float64(len(in)) + 0.5))
	ret.OutputVoltage = int64(math.Floor(float64(ret.OutputVoltage)/float64(len(in)) + 0.5))
	ret.UtilityVoltage = int64(math.Floor(float64(ret.UtilityVoltage)/float64(len(in)) + 0.5))
	ret.RemainingRuntime = int64(math.Floor(float64(ret.RemainingRuntime)/float64(len(in)) + 0.5))

	ret.Status = in[0].Status
	ret.LastTestResult = in[0].LastTestResult
	ret.LastPowerEvent = in[0].LastPowerEvent
	ret.FirmwareVersion = in[0].FirmwareVersion
	ret.ModelName = in[0].ModelName
	ret.Timestamp = in[0].Timestamp
	ret.LineInteraction = in[0].LineInteraction
	ret.AverageOf = len(in)

	return ret, len(in)
}
