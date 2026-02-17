package zenoh

import (
	"testing"
)

func TestNewEncoding(t *testing.T) {
	enc := NewEncoding("test/prefix")
	if enc == nil {
		t.Fatal("NewEncoding() returned nil")
	}
	if enc.Prefix() != "test/prefix" {
		t.Errorf("NewEncoding().Prefix() = %v, want test/prefix", enc.Prefix())
	}
	if enc.Suffix() != "" {
		t.Errorf("NewEncoding().Suffix() = %v, want empty", enc.Suffix())
	}
}

func TestNewEncodingWithSuffix(t *testing.T) {
	enc := NewEncoding("application").WithSuffix("json")
	if enc.Prefix() != "application" {
		t.Errorf("WithSuffix().Prefix() = %v, want application", enc.Prefix())
	}
	if enc.Suffix() != "json" {
		t.Errorf("WithSuffix().Suffix() = %v, want json", enc.Suffix())
	}
}

func TestEncodingString(t *testing.T) {
	tests := []struct {
		name    string
		enc     *Encoding
		wantStr string
	}{
		{
			name:    "prefix only",
			enc:     NewEncoding("text/plain"),
			wantStr: "text/plain",
		},
		{
			name:    "with suffix",
			enc:     NewEncoding("application").WithSuffix("json"),
			wantStr: "application+json",
		},
		{
			name:    "nil",
			enc:     nil,
			wantStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.String(); got != tt.wantStr {
				t.Errorf("Encoding.String() = %v, want %v", got, tt.wantStr)
			}
		})
	}
}

func TestEncodingIsValid(t *testing.T) {
	tests := []struct {
		name    string
		enc     *Encoding
		wantVal bool
	}{
		{
			name:    "valid",
			enc:     NewEncoding("text/plain"),
			wantVal: true,
		},
		{
			name:    "nil",
			enc:     nil,
			wantVal: false,
		},
		{
			name:    "empty prefix",
			enc:     &Encoding{prefix: "", suffix: ""},
			wantVal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.IsValid(); got != tt.wantVal {
				t.Errorf("Encoding.IsValid() = %v, want %v", got, tt.wantVal)
			}
		})
	}
}

func TestEncodingPrefixes(t *testing.T) {
	tests := []struct {
		name string
		enc  *Encoding
		want string
	}{
		{"zenoh/serialized", EncodingZenohSerialized, "zenoh/serialized"},
		{"application/octet-stream", EncodingApplicationOctetStream, "application/octet-stream"},
		{"text/plain", EncodingTextPlain, "text/plain"},
		{"application/json", EncodingApplicationJson, "application/json"},
		{"application/xml", EncodingApplicationXml, "application/xml"},
		{"application/yaml", EncodingApplicationYaml, "application/yaml"},
		{"application/toml", EncodingApplicationToml, "application/toml"},
		{"application/protobuf", EncodingApplicationProtobuf, "application/protobuf"},
		{"application/msgpack", EncodingApplicationMsgPack, "application/msgpack"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.Prefix(); got != tt.want {
				t.Errorf("%s.Prefix() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestEncodingConstructorFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func() *Encoding
		want string
	}{
		{"ZenohSerialized", ZenohSerialized, "zenoh/serialized"},
		{"ApplicationOctetStream", ApplicationOctetStream, "application/octet-stream"},
		{"TextPlain", TextPlain, "text/plain"},
		{"ApplicationJson", ApplicationJson, "application/json"},
		{"ApplicationXml", ApplicationXml, "application/xml"},
		{"ApplicationYaml", ApplicationYaml, "application/yaml"},
		{"ApplicationToml", ApplicationToml, "application/toml"},
		{"ApplicationProtobuf", ApplicationProtobuf, "application/protobuf"},
		{"ApplicationMsgPack", ApplicationMsgPack, "application/msgpack"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.fn()
			if enc.Prefix() != tt.want {
				t.Errorf("%s() = %v, want %v", tt.name, enc.Prefix(), tt.want)
			}
		})
	}
}

func TestEncodingFromStr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		wantP string
		wantS string
	}{
		{
			name:  "prefix only",
			input: "text/plain",
			wantP: "text/plain",
			wantS: "",
		},
		{
			name:  "with suffix",
			input: "application+json",
			wantP: "application",
			wantS: "json",
		},
		{
			name:  "empty",
			input: "",
			wantP: "",
			wantS: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := EncodingFromStr(tt.input)
			if enc == nil && tt.input != "" {
				t.Error("EncodingFromStr() returned nil for non-empty input")
				return
			}
			if enc != nil && enc.Prefix() != tt.wantP {
				t.Errorf("EncodingFromStr().Prefix() = %v, want %v", enc.Prefix(), tt.wantP)
			}
			if enc != nil && enc.Suffix() != tt.wantS {
				t.Errorf("EncodingFromStr().Suffix() = %v, want %v", enc.Suffix(), tt.wantS)
			}
		})
	}
}

func TestEncodingMatches(t *testing.T) {
	tests := []struct {
		name string
		a    *Encoding
		b    *Encoding
		want bool
	}{
		{
			name: "same",
			a:    NewEncoding("text/plain"),
			b:    NewEncoding("text/plain"),
			want: true,
		},
		{
			name: "different",
			a:    NewEncoding("text/plain"),
			b:    NewEncoding("application/json"),
			want: false,
		},
		{
			name: "nil a",
			a:    nil,
			b:    NewEncoding("text/plain"),
			want: false,
		},
		{
			name: "nil b",
			a:    NewEncoding("text/plain"),
			b:    nil,
			want: false,
		},
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Matches(tt.b); got != tt.want {
				t.Errorf("Encoding.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodingEquals(t *testing.T) {
	tests := []struct {
		name string
		a    *Encoding
		b    *Encoding
		want bool
	}{
		{
			name: "same",
			a:    NewEncoding("text/plain"),
			b:    NewEncoding("text/plain"),
			want: true,
		},
		{
			name: "different prefix",
			a:    NewEncoding("text/plain"),
			b:    NewEncoding("application/json"),
			want: false,
		},
		{
			name: "different suffix",
			a:    NewEncoding("application").WithSuffix("json"),
			b:    NewEncoding("application").WithSuffix("xml"),
			want: false,
		},
		{
			name: "nil both",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "nil one",
			a:    nil,
			b:    NewEncoding("text/plain"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equals(tt.b); got != tt.want {
				t.Errorf("Encoding.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodingIsText(t *testing.T) {
	tests := []struct {
		name string
		enc  *Encoding
		want bool
	}{
		{"text/plain", EncodingTextPlain, true},
		{"text/html", NewEncoding("text/html"), true},
		{"text/json", NewEncoding("text/json"), true},
		{"application/json", EncodingApplicationJson, false},
		{"application/octet-stream", EncodingApplicationOctetStream, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.IsText(); got != tt.want {
				t.Errorf("Encoding.IsText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodingIsJson(t *testing.T) {
	tests := []struct {
		name string
		enc  *Encoding
		want bool
	}{
		{"application/json", EncodingApplicationJson, true},
		{"text/json", NewEncoding("text/json"), true},
		{"application/json5", NewEncoding("application/json5"), true},
		{"text/plain", EncodingTextPlain, false},
		{"application/xml", EncodingApplicationXml, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.IsJson(); got != tt.want {
				t.Errorf("Encoding.IsJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodingIsBinary(t *testing.T) {
	tests := []struct {
		name string
		enc  *Encoding
		want bool
	}{
		{"application/octet-stream", EncodingApplicationOctetStream, true},
		{"application/protobuf", EncodingApplicationProtobuf, true},
		{"application/msgpack", EncodingApplicationMsgPack, true},
		{"text/plain", EncodingTextPlain, false},
		{"application/json", EncodingApplicationJson, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enc.IsBinary(); got != tt.want {
				t.Errorf("Encoding.IsBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodingToBytes(t *testing.T) {
	enc := NewEncoding("text/plain")
	data, err := enc.ToBytes()
	if err != nil {
		t.Fatalf("Encoding.ToBytes() error = %v", err)
	}
	if string(data) != "text/plain" {
		t.Errorf("Encoding.ToBytes() = %v, want text/plain", string(data))
	}
}

func TestEncodingToBytesNil(t *testing.T) {
	var enc *Encoding
	_, err := enc.ToBytes()
	if err == nil {
		t.Error("nil Encoding.ToBytes() should return error")
	}
}

func TestEncodingFromBytes(t *testing.T) {
	data := []byte("application/json")
	enc, err := EncodingFromBytes(data)
	if err != nil {
		t.Fatalf("EncodingFromBytes() error = %v", err)
	}
	if enc.Prefix() != "application/json" {
		t.Errorf("EncodingFromBytes().Prefix() = %v, want application/json", enc.Prefix())
	}
}

func TestEncodingFromBytesEmpty(t *testing.T) {
	_, err := EncodingFromBytes([]byte{})
	if err == nil {
		t.Error("EncodingFromBytes(empty) should return error")
	}
}

func TestEncodingRoundTrip(t *testing.T) {
	original := NewEncoding("application").WithSuffix("json")

	data, err := original.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes() error = %v", err)
	}

	decoded, err := EncodingFromBytes(data)
	if err != nil {
		t.Fatalf("EncodingFromBytes() error = %v", err)
	}

	if !original.Equals(decoded) {
		t.Errorf("Round-trip failed: got %v, want %v", decoded.String(), original.String())
	}
}
