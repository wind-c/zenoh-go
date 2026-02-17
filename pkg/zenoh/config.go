package zenoh

import (
	"fmt"
	"unsafe"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

// NewDefaultConfig creates a new zenoh configuration with default settings.
func NewDefaultConfig() (*OwnedConfig, error) {
	cfg, err := cgo.ConfigDefault()
	if err != nil {
		return nil, err
	}
	return &OwnedConfig{ptr: cfg.Ptr, owned: cfg.OwnedPtr()}, nil
}

// ConfigFromFile loads a zenoh configuration from a JSON5 file.
// This is equivalent to zc_config_from_file() in zenoh-c.
func ConfigFromFile(path string) (*OwnedConfig, error) {
	cfg, err := cgo.ConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	return &OwnedConfig{ptr: cfg.Ptr, owned: cfg.OwnedPtr()}, nil
}

// ConfigFromStr loads a zenoh configuration from a JSON5 string.
// This is equivalent to zc_config_from_str() in zenoh-c.
func ConfigFromStr(s string) (*OwnedConfig, error) {
	cfg, err := cgo.ConfigFromStr(s)
	if err != nil {
		return nil, err
	}
	return &OwnedConfig{ptr: cfg.Ptr, owned: cfg.OwnedPtr()}, nil
}

// InsertJSON5 inserts a JSON5 value into the configuration.
// This is equivalent to zc_config_insert_json5() in zenoh-c.
func (c *OwnedConfig) InsertJSON5(key, value string) error {
	if c == nil || !c.IsValid() {
		return ErrInvalidValue
	}
	cfg := cgo.Config{Ptr: c.ptr}
	cfg.SetLoaned(unsafe.Pointer(c.ptr))
	cfg.SetOwned(c.owned)
	return cfg.InsertJSON5(key, value)
}

// EnableQUIC enables QUIC transport on the given listen port.
func (c *OwnedConfig) EnableQUIC(listenPort int) error {
	if c == nil || !c.IsValid() {
		return ErrInvalidValue
	}
	if err := c.InsertJSON5("transport/unicast/quic/enabled", "true"); err != nil {
		return err
	}
	return c.InsertJSON5("transport/unicast/quic/listen", formatPort(listenPort))
}

// EnableQUICClient enables QUIC transport for client mode (no listening).
func (c *OwnedConfig) EnableQUICClient() error {
	if c == nil || !c.IsValid() {
		return ErrInvalidValue
	}
	return c.InsertJSON5("transport/unicast/quic/enabled", "true")
}

// formatPort converts an integer port to a JSON string.
func formatPort(port int) string {
	return fmt.Sprintf("%d", port)
}
