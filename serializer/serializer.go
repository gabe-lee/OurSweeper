package serializer

import (
	"encoding/binary"
	"io"
)

type CodedSerializer interface {
	// Write a unique type code to the Writer, followed by the type data
	CodedSerialize(order binary.ByteOrder) (data []byte, err error)
}

type Serializer interface {
	// Write the type to the Writer
	Serialize(w io.Writer, order binary.ByteOrder) error
}

type Deserializer interface {
	// Read the type from the Reader
	Deserialize(r io.Reader, order binary.ByteOrder) error
}
