package data_buffer

import "time"

const (
	getBufferTimeout = time.Microsecond * 100
	defaultBufferCap = 128
)

type WriteBufferPool struct {
	pool []WriteBuffer
	lock Lock
}

func NewStringBufferPool(initialPoolCap int, initialBufferCap int) WriteBufferPool {
	s := WriteBufferPool{
		pool: make([]WriteBuffer, initialPoolCap),
	}
	for i := range initialPoolCap {
		s.pool[i] = NewWriteBuffer(initialBufferCap)
	}
	return s
}

func (s *WriteBufferPool) GetBuffer() WriteBuffer {
	var buf WriteBuffer
	gotLock := s.lock.TryLockTimeout(getBufferTimeout)
	if gotLock {
		count := len(s.pool)
		if count == 0 {
			buf = NewWriteBuffer(defaultBufferCap)
		} else {
			last := count - 1
			buf = s.pool[last]
			s.pool = s.pool[:last]
		}
		s.lock.Unlock()
	} else {
		buf = NewWriteBuffer(defaultBufferCap)
	}
	return buf
}

func (s *WriteBufferPool) ReleaseBuffer(buf WriteBuffer) {
	buf.Reset()
	go func(buf WriteBuffer) {
		s.lock.Lock()
		s.pool = append(s.pool, buf)
		s.lock.Unlock()
	}(buf)
}

func (s *WriteBufferPool) QuickWriteString(vals ...any) string {
	buf := s.GetBuffer()
	defer s.ReleaseBuffer(buf)
	buf.WriteAny(vals...)
	return buf.StringCopy()
}

func (s *WriteBufferPool) QuickWriteBytes(vals ...any) []byte {
	buf := s.GetBuffer()
	defer s.ReleaseBuffer(buf)
	buf.WriteAny(vals...)
	return buf.BytesCopy()
}

func (s *WriteBufferPool) QuickWriteBuffer(vals ...any) WriteBuffer {
	buf := s.GetBuffer()
	buf.WriteAny(vals...)
	return buf
}
