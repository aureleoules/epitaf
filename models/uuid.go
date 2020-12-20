package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
)

// UUID type
type UUID []byte

// FromUUID returns a UUID struct from a string
func FromUUID(uuid string) (UUID, error) {
	return hex.DecodeString(strings.Replace(uuid, "-", "", -1))
}

// NewUUID util
func NewUUID() UUID {
	id, _ := uuid.NewV4()
	return UUID(id[:])
}

// MarshalJSON interface method
func (id UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON interface method
func (id *UUID) UnmarshalJSON(b []byte) error {
	var err error
	*id, err = FromUUID(strings.Replace(string(b), "\"", "", -1))
	return err
}

func (id UUID) String() string {
	bytes := []byte(id)
	if len(bytes) < 16 {
		return ""
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
}
