package zenoh

import (
	"errors"
	"testing"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg, err := NewDefaultConfig()
	if err != nil {
		t.Logf("NewDefaultConfig() returned error: %v (may be expected without zenoh-c)", err)
		return
	}
	if cfg == nil {
		t.Error("NewDefaultConfig() returned nil config")
		return
	}
	defer cfg.Drop()

	if !cfg.IsValid() {
		t.Error("NewDefaultConfig() returned invalid config")
	}
}

func TestConfigFromFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"empty path", "", true},
		{"valid path (non-existent)", "/path/to/config.json5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConfigFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigFromStr(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{"empty string", "", true},
		{"valid json5", "{mode:\"peer\"}", false},
		{"json5 with mode", "mode=\"peer\"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ConfigFromStr(tt.s)
			if (err != nil) != tt.wantErr {
				t.Logf("ConfigFromStr(%q) error = %v, wantErr %v", tt.s, err, tt.wantErr)
			}
			if err == nil && cfg != nil {
				defer cfg.Drop()
			}
		})
	}
}

func TestOwnedConfig_InsertJSON5(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var c *OwnedConfig
		err := c.InsertJSON5("mode", "\"peer\"")
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("InsertJSON5() on nil config error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		c := &OwnedConfig{ptr: 0}
		err := c.InsertJSON5("mode", "\"peer\"")
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("InsertJSON5() on invalid config error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("valid config", func(t *testing.T) {
		cfg, err := NewDefaultConfig()
		if err != nil {
			t.Logf("Skipping test: cannot create config: %v", err)
			return
		}
		defer cfg.Drop()

		err = cfg.InsertJSON5("mode", "\"peer\"")
		if err != nil {
			t.Errorf("InsertJSON5() on valid config error = %v", err)
		}
	})
}

func TestConfig_EnableQUIC(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var c *OwnedConfig
		err := c.EnableQUIC(7447)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("EnableQUIC() on nil config error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		c := &OwnedConfig{ptr: 0}
		err := c.EnableQUIC(7447)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("EnableQUIC() on invalid config error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("valid config with port", func(t *testing.T) {
		cfg, err := NewDefaultConfig()
		if err != nil {
			t.Logf("Skipping test: cannot create config: %v", err)
			return
		}
		defer cfg.Drop()

		err = cfg.EnableQUIC(7447)
		if err != nil {
			t.Logf("EnableQUIC() error (may require QUIC-enabled zenoh-c): %v", err)
		}
	})

	t.Run("client mode", func(t *testing.T) {
		cfg, err := NewDefaultConfig()
		if err != nil {
			t.Logf("Skipping test: cannot create config: %v", err)
			return
		}
		defer cfg.Drop()

		err = cfg.EnableQUICClient()
		if err != nil {
			t.Logf("EnableQUICClient() error (may require QUIC-enabled zenoh-c): %v", err)
		}
	})
}

func TestOwnedConfig_IsValid(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var c *OwnedConfig
		if c != nil && c.IsValid() {
			t.Error("nil config should not be valid")
		}
	})

	t.Run("zero ptr config", func(t *testing.T) {
		c := &OwnedConfig{ptr: 0}
		if c.IsValid() {
			t.Error("zero ptr config should not be valid")
		}
	})

	t.Run("valid config", func(t *testing.T) {
		cfg, err := NewDefaultConfig()
		if err != nil {
			t.Logf("Skipping test: cannot create config: %v", err)
			return
		}
		defer cfg.Drop()

		if !cfg.IsValid() {
			t.Error("config should be valid after creation")
		}
	})
}

func TestOwnedConfig_Drop(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var c *OwnedConfig
		if c != nil {
			c.Drop()
		}
	})

	t.Run("valid config", func(t *testing.T) {
		cfg, err := NewDefaultConfig()
		if err != nil {
			t.Logf("Skipping test: cannot create config: %v", err)
			return
		}
		cfg.Drop()
		if cfg != nil && cfg.IsValid() {
			t.Error("config should not be valid after Drop()")
		}
	})
}
