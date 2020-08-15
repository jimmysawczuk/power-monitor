package monitor

import (
	"time"
)

type SnapshotSlice []Snapshot
type SnapshotAverageSlice []SnapshotAverage

func (in SnapshotSlice) Len() int {
	return len(in)
}

func (in SnapshotSlice) Rollup(interval time.Duration) SnapshotAverageSlice {
	out := make(SnapshotAverageSlice, 0)

	if len(in) == 0 {
		return out
	}

	lastIndex := 0
	nextCutoff := time.Now().Truncate(interval)

	for i := 0; i < len(in); i++ {
		if in[i].Timestamp.UnixNano() <= nextCutoff.UnixNano() {
			res, len := SnapshotSlice(in[lastIndex:i]).Average()
			res.Timestamp = nextCutoff.Add(interval)
			res.Interval = interval
			if len > 0 {
				out = append(out, res)
			}
			lastIndex = i
			nextCutoff = nextCutoff.Add(-interval)
		}
	}

	res, len := SnapshotSlice(in[lastIndex:]).Average()
	res.Timestamp = nextCutoff.Add(interval)
	res.Interval = interval
	if len > 0 {
		out = append(out, res)
	}

	return out
}

func (in SnapshotSlice) Average() (SnapshotAverage, int) {
	if len(in) == 0 {
		return SnapshotAverage{}, 0
	}

	ret := SnapshotAverage{
		Timestamps: make([]time.Time, 0),
	}

	for i := 0; i < len(in); i++ {
		v := in[i]

		ret.BatteryRemaining += float64(v.BatteryRemaining)
		ret.Load += float64(v.Load)
		ret.BatteryCapacity += float64(v.BatteryCapacity)
		ret.OutputVoltage += float64(v.OutputVoltage)
		ret.UtilityVoltage += float64(v.UtilityVoltage)
		ret.RemainingRuntime += float64(v.RemainingRuntime)

		ret.Timestamps = append(ret.Timestamps, v.Timestamp)
	}

	ret.BatteryRemaining = ret.BatteryRemaining / float64(len(in))
	ret.Load = ret.Load / float64(len(in))
	ret.BatteryCapacity = ret.BatteryCapacity / float64(len(in))
	ret.OutputVoltage = ret.OutputVoltage / float64(len(in))
	ret.UtilityVoltage = ret.UtilityVoltage / float64(len(in))
	ret.RemainingRuntime = ret.RemainingRuntime / float64(len(in))

	ret.AverageOf = len(in)
	return ret, len(in)
}
