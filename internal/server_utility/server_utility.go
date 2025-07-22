package server_utility

import "github.com/gabe-lee/OurSweeper/data_buffer"

type ServerUtilityOptions struct {
	StringBufferPoolCapacity   int
	StringBufferPoolInitBufCap int
}

type ServerUtility struct {
	WriteBufferPool data_buffer.WriteBufferPool
}

func NewServerUtility(options ServerUtilityOptions) ServerUtility {
	return ServerUtility{
		WriteBufferPool: data_buffer.NewStringBufferPool(options.StringBufferPoolCapacity, options.StringBufferPoolInitBufCap),
	}
}
