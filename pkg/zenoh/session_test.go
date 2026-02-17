package zenoh

import (
	"testing"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		config  *OwnedConfig
		wantErr bool
	}{
		{"nil config", nil, true},
		{"invalid config", &OwnedConfig{ptr: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Open(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOwnedSession_Info(t *testing.T) {
	t.Run("nil session", func(t *testing.T) {
		var s *OwnedSession
		_, err := s.Info()
		if err == nil {
			t.Error("Info() on nil session should return error")
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		s := &OwnedSession{ptr: 0}
		_, err := s.Info()
		if err == nil {
			t.Error("Info() on invalid session should return error")
		}
	})
}
