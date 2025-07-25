package server_utility

import "github.com/gabe-lee/OurSweeper/data_buffer"

type (
	WriteBufferPool = data_buffer.WriteBufferPool
)

type ServerUtility struct {
	WriteBufferPool *WriteBufferPool
}

func NewServerUtility() ServerUtility {
	return ServerUtility{
		WriteBufferPool: data_buffer.NewWriteBufferPool(),
	}
}
