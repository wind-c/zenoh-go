package zenoh

import (
	"errors"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

var ErrInvalidEncoding = errors.New("invalid encoding")

type Encoding struct {
	prefix string
	suffix string
}

func NewEncoding(prefix string) *Encoding {
	return &Encoding{prefix: prefix, suffix: ""}
}

func (e *Encoding) Prefix() string {
	if e == nil {
		return ""
	}
	return e.prefix
}

func (e *Encoding) Suffix() string {
	if e == nil {
		return ""
	}
	return e.suffix
}

func (e *Encoding) String() string {
	if e == nil {
		return ""
	}
	if e.suffix != "" {
		return e.prefix + "+" + e.suffix
	}
	return e.prefix
}

func (e *Encoding) IsValid() bool {
	return e != nil && e.prefix != ""
}

func (e *Encoding) WithSuffix(suffix string) *Encoding {
	if e == nil {
		return nil
	}
	return &Encoding{prefix: e.prefix, suffix: suffix}
}

func (e *Encoding) ToBytes() ([]byte, error) {
	if !e.IsValid() {
		return nil, ErrInvalidEncoding
	}
	return []byte(e.String()), nil
}

func EncodingFromBytes(data []byte) (*Encoding, error) {
	if len(data) == 0 {
		return nil, ErrInvalidEncoding
	}
	str := string(data)
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == '+' {
			return &Encoding{prefix: str[:i], suffix: str[i+1:]}, nil
		}
	}
	return &Encoding{prefix: str, suffix: ""}, nil
}

// Predefined encoding prefixes matching zenoh-c conventions.
const (
	EncodingPrefixZenohSerialized        = "zenoh/serialized"
	EncodingPrefixApplicationOctetStream = "application/octet-stream"
	EncodingPrefixTextPlain              = "text/plain"
	EncodingPrefixApplicationJson        = "application/json"
	EncodingPrefixApplicationXml         = "application/xml"
	EncodingPrefixApplicationYaml        = "application/yaml"
	EncodingPrefixApplicationToml        = "application/toml"
	EncodingPrefixApplicationProtobuf    = "application/protobuf"
	EncodingPrefixApplicationMsgPack     = "application/msgpack"
	EncodingPrefixApplicationJson5       = "application/json5"
	EncodingPrefixTextJson               = "text/json"
	EncodingPrefixTextXml                = "text/xml"
	EncodingPrefixTextHtml               = "text/html"
	EncodingPrefixTextCss                = "text/css"
	EncodingPrefixTextJavascript         = "text/javascript"
	EncodingPrefixImageJpeg              = "image/jpeg"
	EncodingPrefixImagePng               = "image/png"
	EncodingPrefixImageGif               = "image/gif"
	EncodingPrefixVideoMp4               = "video/mp4"
	EncodingPrefixVideoWebm              = "video/webm"
	EncodingPrefixAudioMp3               = "audio/mp3"
	EncodingPrefixAudioWav               = "audio/wav"
	EncodingPrefixMultipartMixed         = "multipart/mixed"
	EncodingPrefixMultipartFormData      = "multipart/form-data"
)

// Common encoding presets.
var (
	EncodingZenohSerialized        = NewEncoding(EncodingPrefixZenohSerialized)
	EncodingApplicationOctetStream = NewEncoding(EncodingPrefixApplicationOctetStream)
	EncodingTextPlain              = NewEncoding(EncodingPrefixTextPlain)
	EncodingApplicationJson        = NewEncoding(EncodingPrefixApplicationJson)
	EncodingApplicationXml         = NewEncoding(EncodingPrefixApplicationXml)
	EncodingApplicationYaml        = NewEncoding(EncodingPrefixApplicationYaml)
	EncodingApplicationToml        = NewEncoding(EncodingPrefixApplicationToml)
	EncodingApplicationProtobuf    = NewEncoding(EncodingPrefixApplicationProtobuf)
	EncodingApplicationMsgPack     = NewEncoding(EncodingPrefixApplicationMsgPack)
	EncodingApplicationJson5       = NewEncoding(EncodingPrefixApplicationJson5)
	EncodingTextJson               = NewEncoding(EncodingPrefixTextJson)
	EncodingTextXml                = NewEncoding(EncodingPrefixTextXml)
	EncodingTextHtml               = NewEncoding(EncodingPrefixTextHtml)
	EncodingTextCss                = NewEncoding(EncodingPrefixTextCss)
	EncodingTextJavascript         = NewEncoding(EncodingPrefixTextJavascript)
	EncodingImageJpeg              = NewEncoding(EncodingPrefixImageJpeg)
	EncodingImagePng               = NewEncoding(EncodingPrefixImagePng)
	EncodingImageGif               = NewEncoding(EncodingPrefixImageGif)
	EncodingVideoMp4               = NewEncoding(EncodingPrefixVideoMp4)
	EncodingVideoWebm              = NewEncoding(EncodingPrefixVideoWebm)
	EncodingAudioMp3               = NewEncoding(EncodingPrefixAudioMp3)
	EncodingAudioWav               = NewEncoding(EncodingPrefixAudioWav)
)

// =============================================================================
// Encoding Presets - Constructor Functions
// =============================================================================

// ZenohSerialized returns the zenoh serialized encoding preset.
// This is used for zenoh's internal serialization format.
func ZenohSerialized() *Encoding {
	return EncodingZenohSerialized
}

// ApplicationOctetStream returns the application/octet-stream encoding preset.
// This is used for raw binary data.
func ApplicationOctetStream() *Encoding {
	return EncodingApplicationOctetStream
}

// TextPlain returns the text/plain encoding preset.
// This is used for plain text data.
func TextPlain() *Encoding {
	return EncodingTextPlain
}

// ApplicationJson returns the application/json encoding preset.
// This is used for JSON formatted data.
func ApplicationJson() *Encoding {
	return EncodingApplicationJson
}

// ApplicationXml returns the application/xml encoding preset.
// This is used for XML formatted data.
func ApplicationXml() *Encoding {
	return EncodingApplicationXml
}

// ApplicationYaml returns the application/yaml encoding preset.
// This is used for YAML formatted data.
func ApplicationYaml() *Encoding {
	return EncodingApplicationYaml
}

// ApplicationToml returns the application/toml encoding preset.
// This is used for TOML formatted data.
func ApplicationToml() *Encoding {
	return EncodingApplicationToml
}

// ApplicationProtobuf returns the application/protobuf encoding preset.
// This is used for Protocol Buffers serialized data.
func ApplicationProtobuf() *Encoding {
	return EncodingApplicationProtobuf
}

// ApplicationMsgPack returns the application/msgpack encoding preset.
// This is used for MessagePack serialized data.
func ApplicationMsgPack() *Encoding {
	return EncodingApplicationMsgPack
}

// =============================================================================
// Encoding Resolution
// =============================================================================

// ResolveEncoding resolves an encoding from a string representation.
// It supports both prefix-only and prefix+suffix formats.
func ResolveEncoding(s string) (*Encoding, error) {
	if s == "" {
		return nil, ErrInvalidEncoding
	}
	return EncodingFromStr(s), nil
}

// EncodingFromStr creates an encoding from a string.
func EncodingFromStr(s string) *Encoding {
	if s == "" {
		return nil
	}
	// Check for suffix format: prefix+suffix
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '+' {
			return &Encoding{prefix: s[:i], suffix: s[i+1:]}
		}
	}
	return &Encoding{prefix: s, suffix: ""}
}

// =============================================================================
// Encoding Matching
// =============================================================================

// Matches returns true if this encoding matches another.
// Two encodings match if they have the same prefix.
func (e *Encoding) Matches(other *Encoding) bool {
	if e == nil || other == nil {
		return false
	}
	return e.prefix == other.prefix
}

// Equals returns true if this encoding equals another exactly.
func (e *Encoding) Equals(other *Encoding) bool {
	if e == nil && other == nil {
		return true
	}
	if e == nil || other == nil {
		return false
	}
	return e.prefix == other.prefix && e.suffix == other.suffix
}

// IsText returns true if the encoding represents text data.
func (e *Encoding) IsText() bool {
	if e == nil {
		return false
	}
	return len(e.prefix) >= 4 && e.prefix[:4] == "text"
}

// IsJson returns true if the encoding represents JSON data.
func (e *Encoding) IsJson() bool {
	if e == nil {
		return false
	}
	return e.prefix == "application/json" || e.prefix == "text/json" || e.prefix == "application/json5"
}

// IsBinary returns true if the encoding represents binary data.
func (e *Encoding) IsBinary() bool {
	if e == nil {
		return false
	}
	return e.prefix == "application/octet-stream" ||
		e.prefix == "application/protobuf" ||
		e.prefix == "application/msgpack"
}

// =============================================================================
// CGO Support
// =============================================================================

// toCGO converts the encoding to CGO representation.
// Returns nil as cgo encoding support is not yet implemented.
func (e *Encoding) toCGO() *cgo.Encoding {
	return nil
}

// fromCGO converts CGO representation to Encoding.
func fromCGO(enc *cgo.Encoding) *Encoding {
	return nil
}

// Drop releases the encoding resources.
func (e *Encoding) Drop() error {
	return nil
}
