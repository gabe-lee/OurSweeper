package lockset

type LockSet struct {
	Locks uint64
	XMin  uint8
	XMax  uint8
	YMin  uint8
	YMax  uint8
}

func New() LockSet {
	return LockSet{}
}

func (l *LockSet) AddLock(idx uint, x uint, y uint) {
	val := uint64(1) << idx
	l.Locks = l.Locks | val
	l.XMax = max(uint8(x), l.XMax)
	l.YMax = max(uint8(y), l.YMax)
	l.XMin = min(uint8(x), l.XMin)
	l.YMin = min(uint8(y), l.YMin)
}

func (l LockSet) AlreadyLocked(idx uint) bool {
	val := uint64(1) << idx
	return l.Locks&val > 0
}
