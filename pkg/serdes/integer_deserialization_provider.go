package serdes

import (
	"encoding/binary"
	"fmt"
)

type IntegerDeserializationProvider struct{}

func (IntegerDeserializationProvider) LoadSchema(_ string, _ map[string]string) error {
	return nil
}

func (IntegerDeserializationProvider) Deserialize(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	return fmt.Sprintf("%d", binary.LittleEndian.Uint32(data)), nil
}
