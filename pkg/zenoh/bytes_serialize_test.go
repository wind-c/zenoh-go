package zenoh

import (
	"encoding/base64"
	"testing"
)

func TestBytesSerialize(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "hello",
			input:   []byte("hello"),
			wantErr: false,
		},
		{
			name:    "empty",
			input:   []byte{},
			wantErr: false,
		},
		{
			name:    "binary",
			input:   []byte{0x00, 0x01, 0x02, 0xFF},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBytes(tt.input)
			data, err := b.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			// Check format: [4 bytes: length][data]
			if len(data) < 4 {
				t.Errorf("Serialized data too short: %v", len(data))
			}
		})
	}
}

func TestDeserializeBytes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		wantStr string
	}{
		{
			name:    "hello",
			data:    []byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'},
			wantErr: false,
			wantStr: "hello",
		},
		{
			name:    "too short",
			data:    []byte{0x00, 0x00, 0x01},
			wantErr: true,
		},
		{
			name:    "length mismatch",
			data:    []byte{0x00, 0x00, 0x00, 0x10, 'a'},
			wantErr: true,
		},
		{
			name:    "empty",
			data:    []byte{0x00, 0x00, 0x00, 0x00},
			wantErr: false,
			wantStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := DeserializeBytes(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if b.String() != tt.wantStr {
				t.Errorf("DeserializeBytes().String() = %v, want %v", b.String(), tt.wantStr)
			}
		})
	}
}

func TestSerializeDeserialize(t *testing.T) {
	original := []byte("test data with unicode: 你好")
	b := NewBytes(original)

	serialized, err := b.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	deserialized, err := DeserializeBytes(serialized)
	if err != nil {
		t.Fatalf("DeserializeBytes() error = %v", err)
	}

	if deserialized.String() != string(original) {
		t.Errorf("Round-trip failed: got %v, want %v", deserialized.String(), string(original))
	}
}

func TestBytesToBase64(t *testing.T) {
	b := NewBytes([]byte("hello"))
	encoded, err := b.ToBase64()
	if err != nil {
		t.Fatalf("ToBase64() error = %v", err)
	}
	expected := base64.StdEncoding.EncodeToString([]byte("hello"))
	if encoded != expected {
		t.Errorf("ToBase64() = %v, want %v", encoded, expected)
	}
}

func TestFromBase64(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("hello"))
	b, err := FromBase64(encoded)
	if err != nil {
		t.Fatalf("FromBase64() error = %v", err)
	}
	if b.String() != "hello" {
		t.Errorf("FromBase64().String() = %v, want hello", b.String())
	}
}

func TestFromBase64Invalid(t *testing.T) {
	_, err := FromBase64("!!!invalid!!!")
	if err == nil {
		t.Error("FromBase64() should error on invalid input")
	}
}

func TestBase64RoundTrip(t *testing.T) {
	original := []byte("test data with unicode: 你好")
	b := NewBytes(original)

	encoded, err := b.ToBase64()
	if err != nil {
		t.Fatalf("ToBase64() error = %v", err)
	}

	decoded, err := FromBase64(encoded)
	if err != nil {
		t.Fatalf("FromBase64() error = %v", err)
	}

	if decoded.String() != string(original) {
		t.Errorf("Round-trip failed: got %v, want %v", decoded.String(), string(original))
	}
}
