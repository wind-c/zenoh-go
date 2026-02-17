package zenoh

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"
)

var ErrInvalidBytes = errors.New("invalid zenoh bytes")

// =============================================================================
// OwnedBytes - Owned Bytes Type
// =============================================================================

type OwnedBytes struct {
	ptr  uintptr
	data []byte
}

func NewOwnedBytes() (*OwnedBytes, error) {
	return &OwnedBytes{
		ptr:  0,
		data: make([]byte, 0),
	}, nil
}

func NewOwnedBytesFromSlice(data []byte) (*OwnedBytes, error) {
	if data == nil {
		return nil, ErrInvalidBytes
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return &OwnedBytes{
		ptr:  0,
		data: cp,
	}, nil
}

func NewOwnedBytesFromString(s string) (*OwnedBytes, error) {
	return NewOwnedBytesFromSlice([]byte(s))
}

func FromOwnedBytes(owned *OwnedBytes) (*Bytes, error) {
	if owned == nil {
		return nil, ErrInvalidBytes
	}
	if !owned.IsValid() {
		return nil, ErrInvalidBytes
	}
	if owned.ptr == 0 && (owned.data == nil || len(owned.data) == 0) {
		return &Bytes{data: make([]byte, 0)}, nil
	}
	if owned.ptr == 0 && len(owned.data) > 0 {
		data := make([]byte, len(owned.data))
		copy(data, owned.data)
		return &Bytes{data: data}, nil
	}
	data := make([]byte, len(owned.data))
	copy(data, owned.data)
	return &Bytes{data: data}, nil
}

func (b *OwnedBytes) Drop() error {
	if b.ptr == 0 && len(b.data) == 0 {
		return nil
	}
	log.Print("[zenoh] bytes dropped")
	b.ptr = 0
	b.data = nil
	return nil
}

func (b *OwnedBytes) IsValid() bool {
	return (b.ptr != 0) || (b.data != nil)
}

func (b *OwnedBytes) Data() []byte {
	if b == nil {
		return nil
	}
	return b.data
}

func (b *OwnedBytes) Len() int {
	return len(b.Data())
}

func (b *OwnedBytes) String() string {
	return string(b.Data())
}

// =============================================================================
// Bytes - Loaned Bytes Type
// =============================================================================

type Bytes struct {
	ptr  uintptr
	data []byte
}

func NewBytes(data []byte) *Bytes {
	if data == nil {
		return &Bytes{data: make([]byte, 0)}
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return &Bytes{data: cp}
}

func NewBytesFromString(s string) *Bytes {
	return NewBytes([]byte(s))
}

func (b *Bytes) IsValid() bool {
	return b != nil && (b.ptr != 0 || len(b.data) >= 0)
}

func (b *Bytes) Data() []byte {
	if b == nil {
		return nil
	}
	return b.data
}

func (b *Bytes) Len() int {
	if b == nil {
		return 0
	}
	return len(b.Data())
}

func (b *Bytes) String() string {
	if b == nil {
		return ""
	}
	return string(b.Data())
}

// =============================================================================
// Serialization/Deserialization
// =============================================================================

func (b *Bytes) Serialize() ([]byte, error) {
	if !b.IsValid() {
		return nil, ErrInvalidBytes
	}
	data := b.Data()
	if data == nil {
		return nil, ErrInvalidBytes
	}
	result := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(result[0:4], uint32(len(data)))
	copy(result[4:], data)
	return result, nil
}

func DeserializeBytes(data []byte) (*Bytes, error) {
	if len(data) < 4 {
		return nil, ErrInvalidBytes
	}
	length := binary.BigEndian.Uint32(data[0:4])
	if int(length) != len(data)-4 {
		return nil, ErrInvalidBytes
	}
	return NewBytes(data[4 : 4+length]), nil
}

func (b *Bytes) ToBase64() (string, error) {
	if !b.IsValid() {
		return "", ErrInvalidBytes
	}
	return base64.StdEncoding.EncodeToString(b.Data()), nil
}

func FromBase64(encoded string) (*Bytes, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return NewBytes(data), nil
}
