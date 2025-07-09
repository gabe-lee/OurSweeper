package wire_serializer

import (
	"encoding/binary"
	"fmt"
	"io"
)

type WireSerializer interface {
	// This should write the specified code first, followed by the internal data
	//
	// If the passed code is invalid for the type, method should panic
	// (The caller is responsible for concretely ensuring codes are only passed
	// to types that can handle them for serialization)
	WriteWire(w io.Writer, order binary.ByteOrder, code uint32) error
	// This should assume the type code has already been read (and is passed as the `code` parameter)
	//
	// If the passed code is invalid for the type, method should panic
	// (The caller is responsible for concretely ensuring codes are only passed
	// to types that can handle them for serialization)
	ReadWire(r io.Reader, order binary.ByteOrder, code uint32) error
}

func MakeError(badCode uint32, typeName string, goodCodes []uint32) error {
	return fmt.Errorf("wire code %d is invalid for type %s, it only accepts codes: %v", badCode, typeName, goodCodes)
}
