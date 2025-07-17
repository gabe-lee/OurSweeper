package serializer

import (
	"encoding/binary"
	"io"
)

type Opt struct {
	ByteOrder  binary.ByteOrder
	IncludeTag bool
	TagValue   int32
}

func Default() Opt {
	return Opt{
		ByteOrder: binary.LittleEndian,
	}
}

func (o Opt) LittleEndian() Opt {
	o.ByteOrder = binary.LittleEndian
	return o
}
func (o Opt) BigEndian() Opt {
	o.ByteOrder = binary.BigEndian
	return o
}
func (o Opt) NativeEndian() Opt {
	o.ByteOrder = binary.NativeEndian
	return o
}
func (o Opt) WithTag(tag int32) Opt {
	o.IncludeTag = true
	o.TagValue = tag
	return o
}

type WireWriter interface {
	// Write the serialized type data to the provided writer
	WireWrite(writer io.Writer) error
}

type WireReader interface {
	// Read the serialized type data from the provided byte slice into the type
	WireRead(reader io.Reader) error
}

type SerialDeserializer interface {
	WireWriter
	WireReader
}

// This interface provides size hints for the serialization process
type SerialHintSizer interface {
	// This should return the minimum size this type can serialize into
	//
	// This should be considered a fast *hint* for preallocating buffer space,
	// and not an exact minimum (true size may be smaller)
	MinSerialSize() int
	// This should return the maximum size this type can serialize into
	//
	// This should be considered a fast *hint* for preallocating buffer space,
	// and not an exact maximim (true size might be larger)
	MaxSerialSize() int
}

type SerialExactSizer interface {
	// This should return the exact byte length required to serialize
	// the type in its current state
	ExactSerialSize() int
}
