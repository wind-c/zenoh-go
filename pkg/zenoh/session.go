package zenoh

import (
	"log"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

// Open opens a zenoh session with the given configuration.
func Open(config *OwnedConfig) (*OwnedSession, error) {
	if config == nil || !config.IsValid() {
		return nil, ErrInvalidValue
	}
	log.Printf("[zenoh-go] Open called, config ptr: %d, IsValid: %v", config.ptr, config.IsValid())
	cfg := cgo.Config{Ptr: uintptr(config.ptr)}
	cfg.SetOwnedPtr(config.owned)
	log.Printf("[zenoh-go] Created cgo.Config with Ptr: %d", cfg.Ptr)
	s, err := cfg.Open()
	if err != nil {
		log.Printf("[zenoh-go] cfg.Open() failed: %v", err)
		return nil, err
	}
	return &OwnedSession{ptr: s.Ptr, owned: s.OwnedPtr()}, nil
}
