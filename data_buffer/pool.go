package data_buffer

import (
	"math/bits"
	"sync"
)

type SizeClass int

const (
	minPoolClass = 6
	maxPoolClass = 13
	poolClasses  = (maxPoolClass + 1) - minPoolClass

	class_64    = 0
	class_128   = 1
	class_256   = 2
	class_512   = 3
	class_1024  = 4
	class_2048  = 5
	class_4096  = 6
	class_Large = 7

	Size_64    = 64
	Size_128   = 128
	Size_256   = 256
	Size_512   = 512
	Size_1024  = 1024
	Size_2048  = 2048
	Size_4096  = 4096
	Size_Large = 4097
)

var (
	Class_64    = SizeClass(class_64)
	Class_128   = SizeClass(class_128)
	Class_256   = SizeClass(class_256)
	Class_512   = SizeClass(class_512)
	Class_1024  = SizeClass(class_1024)
	Class_2048  = SizeClass(class_2048)
	Class_4096  = SizeClass(class_4096)
	Class_Large = SizeClass(class_Large)
)

type WriteBufferPool struct {
	noCopy

	pools [poolClasses]sync.Pool
}

func NewWriteBufferPool() *WriteBufferPool {
	p := WriteBufferPool{}
	pp := &p
	pp.pools[class_64] = sync.Pool{}
	pp.pools[class_64].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_64),
		}
	}
	pp.pools[class_128] = sync.Pool{}
	pp.pools[class_128].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_128),
		}
	}
	pp.pools[class_256] = sync.Pool{}
	pp.pools[class_256].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_256),
		}
	}
	pp.pools[class_512] = sync.Pool{}
	pp.pools[class_512].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_512),
		}
	}
	pp.pools[class_1024] = sync.Pool{}
	pp.pools[class_1024].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_1024),
		}
	}
	pp.pools[class_2048] = sync.Pool{}
	pp.pools[class_2048].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_2048),
		}
	}
	pp.pools[class_4096] = sync.Pool{}
	pp.pools[class_4096].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_4096),
		}
	}
	pp.pools[class_Large] = sync.Pool{}
	pp.pools[class_Large].New = func() any {
		return &WriteBuffer{
			data: make([]byte, 0, Size_Large),
		}
	}
	return pp
}

func getSizeClass(size int) SizeClass {
	s := uint(size)
	s -= 1
	s |= s >> 1
	s |= s >> 2
	s |= s >> 4
	s |= s >> 8
	s |= s >> 16
	if bits.UintSize == 64 {
		s |= s >> 32
	}
	s += 1
	class := bits.TrailingZeros(s)
	class = min(maxPoolClass, max(class, minPoolClass))
	class -= minPoolClass
	return SizeClass(class)
}

func (s *WriteBufferPool) GetBuffer(size int) *WriteBuffer {
	class := getSizeClass(size)
	return s.getBufferByClassInternal(class)
}
func (s *WriteBufferPool) GetBufferByClass(class SizeClass) *WriteBuffer {
	class = min(class_Large, max(class_64, class))
	return s.getBufferByClassInternal(class)
}
func (s *WriteBufferPool) getBufferByClassInternal(class SizeClass) *WriteBuffer {
	return s.pools[class].Get().(*WriteBuffer)
}

func (s *WriteBufferPool) ReleaseBuffer(buf *WriteBuffer) {
	buf.Reset()
	class := getSizeClass(cap(buf.data))
	s.pools[class].Put(buf)
}

func (s *WriteBufferPool) QuickWriteTextAndCopyString(sizeClass SizeClass, vals ...any) string {
	class := min(class_Large, max(class_64, sizeClass))
	buf := s.GetBufferByClass(class)
	defer s.ReleaseBuffer(buf)
	buf.WriteText(vals...)
	return buf.StringCopy()
}

func (s *WriteBufferPool) QuickWriteTextAndCopyBytes(sizeClass SizeClass, vals ...any) []byte {
	class := min(class_Large, max(class_64, sizeClass))
	buf := s.GetBufferByClass(class)
	defer s.ReleaseBuffer(buf)
	buf.WriteText(vals...)
	return buf.BytesCopy()
}

func (s *WriteBufferPool) QuickWriteText(sizeClass SizeClass, vals ...any) *WriteBuffer {
	class := min(class_Large, max(class_64, sizeClass))
	buf := s.GetBufferByClass(class)
	buf.WriteText(vals...)
	return buf
}
