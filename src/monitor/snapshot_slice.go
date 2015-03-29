package monitor

type SnapshotSlice []Snapshot

func (in SnapshotSlice) Filter(test func(Snapshot) bool) (out SnapshotSlice) {
	out = make(SnapshotSlice, 0)

	for i := 0; i < len(in); i++ {
		if test(in[i]) {
			out = append(out, in[i])
		}
	}

	return out
}
