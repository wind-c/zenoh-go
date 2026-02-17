package zenoh

import (
	"testing"
)

func TestNewBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantLen  int
		wantStr  string
		wantTrue bool
	}{
		{
			name:     "nil input",
			input:    nil,
			wantLen:  0,
			wantStr:  "",
			wantTrue: true,
		},
		{
			name:     "empty slice",
			input:    []byte{},
			wantLen:  0,
			wantStr:  "",
			wantTrue: true,
		},
		{
			name:     "hello world",
			input:    []byte("hello world"),
			wantLen:  11,
			wantStr:  "hello world",
			wantTrue: true,
		},
		{
			name:     "binary data",
			input:    []byte{0x00, 0x01, 0x02, 0xFF},
			wantLen:  4,
			wantStr:  "\x00\x01\x02\xff",
			wantTrue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBytes(tt.input)
			if got := b.IsValid(); got != tt.wantTrue {
				t.Errorf("NewBytes().IsValid() = %v, want %v", got, tt.wantTrue)
			}
			if got := b.Len(); got != tt.wantLen {
				t.Errorf("NewBytes().Len() = %v, want %v", got, tt.wantLen)
			}
			if got := b.String(); got != tt.wantStr {
				t.Errorf("NewBytes().String() = %v, want %v", got, tt.wantStr)
			}
		})
	}
}

func TestNewBytesFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantStr string
	}{
		{
			name:    "empty string",
			input:   "",
			wantLen: 0,
			wantStr: "",
		},
		{
			name:    "hello",
			input:   "hello",
			wantLen: 5,
			wantStr: "hello",
		},
		{
			name:    "unicode",
			input:   "你好",
			wantLen: 6,
			wantStr: "你好",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBytesFromString(tt.input)
			if got := b.Len(); got != tt.wantLen {
				t.Errorf("NewBytesFromString().Len() = %v, want %v", got, tt.wantLen)
			}
			if got := b.String(); got != tt.wantStr {
				t.Errorf("NewBytesFromString().String() = %v, want %v", got, tt.wantStr)
			}
		})
	}
}

func TestBytesData(t *testing.T) {
	original := []byte("test data")
	b := NewBytes(original)

	// Modify original to verify copy
	original[0] = 'X'

	data := b.Data()
	if string(data) != "test data" {
		t.Errorf("Bytes().Data() = %v, want test data", string(data))
	}
}

func TestBytesNil(t *testing.T) {
	var b *Bytes

	if b.IsValid() {
		t.Error("nil Bytes should not be valid")
	}
	if b.Len() != 0 {
		t.Errorf("nil Bytes.Len() = %v, want 0", b.Len())
	}
	if b.String() != "" {
		t.Errorf("nil Bytes.String() = %v, want empty string", b.String())
	}
	if b.Data() != nil {
		t.Errorf("nil Bytes.Data() = %v, want nil", b.Data())
	}
}

func TestOwnedBytesNewOwnedBytes(t *testing.T) {
	ob, err := NewOwnedBytes()
	if err != nil {
		t.Errorf("NewOwnedBytes() error = %v", err)
	}
	if !ob.IsValid() {
		t.Error("NewOwnedBytes() should be valid (empty is valid)")
	}
	if ob.Len() != 0 {
		t.Errorf("NewOwnedBytes().Len() = %v, want 0", ob.Len())
	}
}

func TestOwnedBytesFromSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
		wantLen int
	}{
		{
			name:    "nil",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "valid data",
			input:   []byte("hello"),
			wantErr: false,
			wantLen: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob, err := NewOwnedBytesFromSlice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOwnedBytesFromSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if ob.Len() != tt.wantLen {
				t.Errorf("NewOwnedBytesFromSlice().Len() = %v, want %v", ob.Len(), tt.wantLen)
			}
		})
	}
}

func TestOwnedBytesFromString(t *testing.T) {
	ob, err := NewOwnedBytesFromString("test string")
	if err != nil {
		t.Errorf("NewOwnedBytesFromString() error = %v", err)
	}
	if ob.String() != "test string" {
		t.Errorf("NewOwnedBytesFromString().String() = %v, want test string", ob.String())
	}
}

func TestOwnedBytesDrop(t *testing.T) {
	ob, err := NewOwnedBytesFromString("test")
	if err != nil {
		t.Fatalf("NewOwnedBytesFromString() error = %v", err)
	}

	if err := ob.Drop(); err != nil {
		t.Errorf("Drop() error = %v", err)
	}

	if ob.IsValid() {
		t.Error("OwnedBytes should be invalid after Drop()")
	}
}

func TestOwnedBytesDropMultipleTimes(t *testing.T) {
	ob, err := NewOwnedBytesFromString("test")
	if err != nil {
		t.Fatalf("NewOwnedBytesFromString() error = %v", err)
	}

	// First drop
	if err := ob.Drop(); err != nil {
		t.Errorf("Drop() error = %v", err)
	}

	// Second drop should be no-op
	if err := ob.Drop(); err != nil {
		t.Errorf("Drop() second call error = %v", err)
	}
}

func TestFromOwnedBytes(t *testing.T) {
	ob, err := NewOwnedBytesFromString("owned data")
	if err != nil {
		t.Fatalf("NewOwnedBytesFromString() error = %v", err)
	}
	defer ob.Drop()

	b, err := FromOwnedBytes(ob)
	if err != nil {
		t.Errorf("FromOwnedBytes() error = %v", err)
	}
	if !b.IsValid() {
		t.Error("FromOwnedBytes() should be valid")
	}
	if b.String() != "owned data" {
		t.Errorf("FromOwnedBytes().String() = %v, want owned data", b.String())
	}
}

func TestFromOwnedBytesNil(t *testing.T) {
	_, err := FromOwnedBytes(nil)
	if err != ErrInvalidBytes {
		t.Errorf("FromOwnedBytes(nil) error = %v, want ErrInvalidBytes", err)
	}
}

func TestFromOwnedBytesInvalid(t *testing.T) {
	ob := &OwnedBytes{}
	_, err := FromOwnedBytes(ob)
	if err != ErrInvalidBytes {
		t.Errorf("FromOwnedBytes(invalid) error = %v, want ErrInvalidBytes", err)
	}
}
