package wire_serializer

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	DEFAULT_CODE uint32 = 0
)

type WireSerializer interface {
	// Serialize the type for transfer over a network or into a file.
	// 	* 'mode' can be any arbitrary uint32 that signals to the type what procedure to use when
	// 	reading/writing the data, usually a value specified within the same package as the type.
	// 	* If there is only ONE way to serialize this type, the implementor can skip writing the mode to the writer,
	//	as long as that is mirrored in the implementation of 'ReadWire'
	// 	* If multiple modes are available, this function MUST write the uint32 'mode' first before any other data.
	WriteWire(w io.Writer, order binary.ByteOrder, mode uint32) error
	// Deserialize the type from a network or file reader.
	// 	* If there is only ONE way to deserialize this type, the implementor can skip reading the 'mode'
	//	uint32, as long as that is mirrored in the implementation of 'WriteWire'
	// 	* If there are multiple ways to deserialize this type, the implementor MUST read the uint32 'mode'
	//	first before any other data
	ReadWire(r io.Reader, order binary.ByteOrder) error
}

func MakeError(badMode uint32, typeName string, goodModes []uint32) error {
	return fmt.Errorf("mode %d is invalid for type %s, it only accepts modes: %v", badMode, typeName, goodModes)
}
