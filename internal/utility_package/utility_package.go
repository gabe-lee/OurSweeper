package utility_package

import "github.com/gabe-lee/OurSweeper/data_buffer"

type (
	WriteBufferPool = data_buffer.WriteBufferPool
)

type UtilityPackage struct {
	WriteBufferPool *WriteBufferPool
}

func NewServerUtility() UtilityPackage {
	return UtilityPackage{
		WriteBufferPool: data_buffer.NewWriteBufferPool(),
	}
}
