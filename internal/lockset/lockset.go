package lockset

type LockSet struct {
	Locks uint64
	XMin  uint8
	XMax  uint8
	YMin  uint8
	YMax  uint8
}

func (l *LockSet) AddLock(idx int, x int, y int) {
	val := uint64(1) << idx
	l.Locks = l.Locks | val
	l.XMax = max(uint8(x), l.XMax)
	l.YMax = max(uint8(y), l.YMax)
	l.XMin = min(uint8(x), l.XMin)
	l.YMin = min(uint8(y), l.YMin)
}

func (l LockSet) AlreadyLocked(idx int) bool {
	val := uint64(1) << idx
	return l.Locks&val > 0
}

type QuadLock struct {
	LockIndexes [4]int
	Len         int
}
