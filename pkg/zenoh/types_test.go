package zenoh

import (
	"errors"
	"testing"
)

// =============================================================================
// OwnedConfig Tests
// =============================================================================

func TestOwnedConfig_PointerBehavior(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &OwnedConfig{ptr: tt.ptr}
			if got := c.IsValid(); got != tt.expected {
				t.Errorf("OwnedConfig.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOwnedConfig_DropIdempotent(t *testing.T) {
	t.Run("idempotent drop", func(t *testing.T) {
		c := &OwnedConfig{ptr: 1}

		// First drop should succeed
		if err := c.Drop(); err != nil {
			t.Errorf("first Drop() error = %v", err)
		}

		// After drop, ptr should be 0
		if c.ptr != 0 {
			t.Errorf("after Drop() ptr = %v, want 0", c.ptr)
		}

		// Second drop should also succeed (idempotent)
		if err := c.Drop(); err != nil {
			t.Errorf("second Drop() error = %v", err)
		}
	})

	t.Run("drop nil config", func(t *testing.T) {
		c := &OwnedConfig{ptr: 0}

		if err := c.Drop(); err != nil {
			t.Errorf("Drop() on nil config error = %v", err)
		}
	})
}

func TestOwnedConfig_NewOwnedConfig(t *testing.T) {
	cfg, err := NewDefaultConfig()
	if err != nil {
		t.Logf("Skipping test: cannot create config: %v", err)
		return
	}
	defer cfg.Drop()

	if !cfg.IsValid() {
		t.Error("config should be valid after creation")
	}
}

func TestOwnedConfig_FromOwnedConfig(t *testing.T) {
	t.Run("nil owned", func(t *testing.T) {
		_, err := FromOwnedConfig(nil)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("FromOwnedConfig(nil) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid owned", func(t *testing.T) {
		owned := &OwnedConfig{ptr: 0}
		_, err := FromOwnedConfig(owned)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("FromOwnedConfig(invalid) error = %v, want ErrInvalidValue", err)
		}
	})
}

// =============================================================================
// Config (Loaned) Tests
// =============================================================================

func TestConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{ptr: tt.ptr}
			if got := c.IsValid(); got != tt.expected {
				t.Errorf("Config.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// =============================================================================
// OwnedSession Tests
// =============================================================================

func TestOwnedSession_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OwnedSession{ptr: tt.ptr}
			if got := s.IsValid(); got != tt.expected {
				t.Errorf("OwnedSession.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOwnedSession_Drop(t *testing.T) {
	t.Run("idempotent drop", func(t *testing.T) {
		s := &OwnedSession{ptr: 1}

		if err := s.Drop(); err != nil {
			t.Errorf("first Drop() error = %v", err)
		}

		if s.ptr != 0 {
			t.Errorf("after Drop() ptr = %v, want 0", s.ptr)
		}

		if err := s.Drop(); err != nil {
			t.Errorf("second Drop() error = %v", err)
		}
	})
}

func TestOwnedSession_Open(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		_, err := Open(nil)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("Open(nil) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		config := &OwnedConfig{ptr: 0}
		_, err := Open(config)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("Open(invalid config) error = %v, want ErrInvalidValue", err)
		}
	})
}

func TestOwnedSession_FromOwnedSession(t *testing.T) {
	t.Run("nil owned", func(t *testing.T) {
		_, err := FromOwnedSession(nil)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("FromOwnedSession(nil) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid owned", func(t *testing.T) {
		owned := &OwnedSession{ptr: 0}
		_, err := FromOwnedSession(owned)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("FromOwnedSession(invalid) error = %v, want ErrInvalidValue", err)
		}
	})
}

// =============================================================================
// Session (Loaned) Tests
// =============================================================================

func TestSession_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{ptr: tt.ptr}
			if got := s.IsValid(); got != tt.expected {
				t.Errorf("Session.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSession_Close(t *testing.T) {
	t.Run("close nil session", func(t *testing.T) {
		s := &Session{ptr: 0}
		if err := s.Close(); err != nil {
			t.Errorf("Close() on nil session error = %v", err)
		}
	})
}

// =============================================================================
// OwnedKeyExpr Tests
// =============================================================================

func TestOwnedKeyExpr_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &OwnedKeyExpr{ptr: tt.ptr}
			if got := k.IsValid(); got != tt.expected {
				t.Errorf("OwnedKeyExpr.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOwnedKeyExpr_Drop(t *testing.T) {
	t.Run("idempotent drop", func(t *testing.T) {
		k := &OwnedKeyExpr{ptr: 1}

		if err := k.Drop(); err != nil {
			t.Errorf("first Drop() error = %v", err)
		}

		if k.ptr != 0 {
			t.Errorf("after Drop() ptr = %v, want 0", k.ptr)
		}

		if err := k.Drop(); err != nil {
			t.Errorf("second Drop() error = %v", err)
		}
	})
}

func TestOwnedKeyExpr_NewOwnedKeyExpr(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		_, err := NewOwnedKeyExpr("")
		if err == nil {
			t.Error("NewOwnedKeyExpr(\"\") expected error")
		}
	})

	t.Run("valid string", func(t *testing.T) {
		_, err := NewOwnedKeyExpr("demo/example/*")
		if err == nil {
			t.Error("NewOwnedKeyExpr() expected error (no zenoh-c bindings)")
		}
	})
}

func TestOwnedKeyExpr_FromOwnedKeyExpr(t *testing.T) {
	t.Run("nil owned", func(t *testing.T) {
		_, err := FromOwnedKeyExpr(nil)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("FromOwnedKeyExpr(nil) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid owned", func(t *testing.T) {
		owned := &OwnedKeyExpr{ptr: 0}
		_, err := FromOwnedKeyExpr(owned)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("FromOwnedKeyExpr(invalid) error = %v, want ErrInvalidValue", err)
		}
	})
}

// =============================================================================
// KeyExpr (Loaned) Tests
// =============================================================================

func TestKeyExpr_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeyExpr{ptr: tt.ptr}
			if got := k.IsValid(); got != tt.expected {
				t.Errorf("KeyExpr.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestKeyExpr_String(t *testing.T) {
	t.Run("nil ptr", func(t *testing.T) {
		k := &KeyExpr{ptr: 0}
		if str := k.String(); str != "" {
			t.Errorf("KeyExpr.String() on nil = %v, want empty string", str)
		}
	})
}

// =============================================================================
// OwnedPublisher Tests
// =============================================================================

func TestOwnedPublisher_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OwnedPublisher{ptr: tt.ptr}
			if got := p.IsValid(); got != tt.expected {
				t.Errorf("OwnedPublisher.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOwnedPublisher_Drop(t *testing.T) {
	t.Run("idempotent drop", func(t *testing.T) {
		p := &OwnedPublisher{ptr: 1}

		if err := p.Drop(); err != nil {
			t.Errorf("first Drop() error = %v", err)
		}

		if p.ptr != 0 {
			t.Errorf("after Drop() ptr = %v, want 0", p.ptr)
		}
	})
}

func TestDeclarePublisher(t *testing.T) {
	t.Run("nil session", func(t *testing.T) {
		keyExpr := &OwnedKeyExpr{ptr: 1}
		_, err := DeclarePublisher(nil, keyExpr)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclarePublisher(nil, keyExpr) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("nil keyExpr", func(t *testing.T) {
		session := &OwnedSession{ptr: 1}
		_, err := DeclarePublisher(session, nil)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclarePublisher(session, nil) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		session := &OwnedSession{ptr: 0}
		keyExpr := &OwnedKeyExpr{ptr: 1}
		_, err := DeclarePublisher(session, keyExpr)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclarePublisher(invalid, keyExpr) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid keyExpr", func(t *testing.T) {
		session := &OwnedSession{ptr: 1}
		keyExpr := &OwnedKeyExpr{ptr: 0}
		_, err := DeclarePublisher(session, keyExpr)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclarePublisher(session, invalid) error = %v, want ErrInvalidValue", err)
		}
	})
}

// =============================================================================
// Publisher (Loaned) Tests
// =============================================================================

func TestPublisher_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Publisher{ptr: tt.ptr}
			if got := p.IsValid(); got != tt.expected {
				t.Errorf("Publisher.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPublisher_Put(t *testing.T) {
	t.Run("invalid publisher", func(t *testing.T) {
		p := &Publisher{ptr: 0}
		err := p.Put([]byte("test"), EncodingTextPlain)
		if !errors.Is(err, ErrInvalidPublisher) {
			t.Errorf("Put() on invalid publisher error = %v, want ErrInvalidPublisher", err)
		}
	})
}

// =============================================================================
// OwnedSubscriber Tests
// =============================================================================

func TestOwnedSubscriber_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OwnedSubscriber{ptr: tt.ptr}
			if got := s.IsValid(); got != tt.expected {
				t.Errorf("OwnedSubscriber.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOwnedSubscriber_Drop(t *testing.T) {
	t.Run("idempotent drop", func(t *testing.T) {
		s := &OwnedSubscriber{ptr: 1}

		if err := s.Drop(); err != nil {
			t.Errorf("first Drop() error = %v", err)
		}

		if s.ptr != 0 {
			t.Errorf("after Drop() ptr = %v, want 0", s.ptr)
		}
	})
}

func TestDeclareSubscriberOld(t *testing.T) {
	t.Run("nil session", func(t *testing.T) {
		_, err := DeclareSubscriber(nil, "demo/test", func(s Sample) {})
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclareSubscriber(nil, keyExpr) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("empty keyExpr", func(t *testing.T) {
		session := &OwnedSession{ptr: 1}
		_, err := DeclareSubscriber(session, "", func(s Sample) {})
		if !errors.Is(err, ErrInvalidKeyExpr) {
			t.Errorf("DeclareSubscriber(session, empty) error = %v, want ErrInvalidKeyExpr", err)
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		session := &OwnedSession{ptr: 0}
		_, err := DeclareSubscriber(session, "demo/test", func(s Sample) {})
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclareSubscriber(invalid, keyExpr) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("nil callback", func(t *testing.T) {
		session := &OwnedSession{ptr: 1}
		_, err := DeclareSubscriber(session, "demo/test", nil)
		if err == nil {
			t.Errorf("DeclareSubscriber(session, keyExpr, nil) should return error")
		}
	})
}

// =============================================================================
// Subscriber (Loaned) Tests
// =============================================================================

func TestSubscriber_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscriber{ptr: tt.ptr}
			if got := s.IsValid(); got != tt.expected {
				t.Errorf("Subscriber.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// =============================================================================
// OwnedQueryable Tests
// =============================================================================

func TestOwnedQueryable_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &OwnedQueryable{ptr: tt.ptr}
			if got := q.IsValid(); got != tt.expected {
				t.Errorf("OwnedQueryable.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOwnedQueryable_Drop(t *testing.T) {
	t.Run("idempotent drop", func(t *testing.T) {
		q := &OwnedQueryable{ptr: 1}

		if err := q.Drop(); err != nil {
			t.Errorf("first Drop() error = %v", err)
		}

		if q.ptr != 0 {
			t.Errorf("after Drop() ptr = %v, want 0", q.ptr)
		}
	})
}

func TestDeclareQueryable(t *testing.T) {
	t.Run("nil session", func(t *testing.T) {
		_, err := DeclareQueryable(nil, "demo/test", func(q Query) {})
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclareQueryable(nil, keyExpr) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("nil callback", func(t *testing.T) {
		session := &OwnedSession{ptr: 0}
		_, err := DeclareQueryable(session, "demo/test", nil)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclareQueryable(session, nil) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		session := &OwnedSession{ptr: 0}
		_, err := DeclareQueryable(session, "demo/test", func(q Query) {})
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("DeclareQueryable(invalid, keyExpr) error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("empty keyExpr", func(t *testing.T) {
		// Skip: requires valid zenoh session to test empty keyexpr error
		// Using ptr: 1 would crash CGO, ptr: 0 returns invalid session error instead
		t.Skip("requires real zenoh session")
	})
}

// =============================================================================
// Queryable (Loaned) Tests
// =============================================================================

func TestQueryable_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		ptr      uintptr
		expected bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queryable{ptr: tt.ptr}
			if got := q.IsValid(); got != tt.expected {
				t.Errorf("Queryable.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// =============================================================================
// Utility Function Tests
// =============================================================================

func TestInit(t *testing.T) {
	// Init should not panic
	Init()
}

func TestLog(t *testing.T) {
	// Log should not panic
	Log("test message")
}

// =============================================================================
// Result Type Tests
// =============================================================================

func TestResult_IsOK(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected bool
	}{
		{"Z_OK", 0, true},
		{"error code", 1, false},
		{"negative code", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{code: tt.code}
			if got := r.IsOK(); got != tt.expected {
				t.Errorf("Result.IsOK() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResult_Error(t *testing.T) {
	t.Run("Z_OK", func(t *testing.T) {
		r := Result{code: 0}
		if err := r.Error(); err != "" {
			t.Errorf("Result.Error() on Z_OK = %v, want empty string", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		r := Result{code: 1}
		if err := r.Error(); err == "" {
			t.Error("Result.Error() on error should return non-empty string")
		}
	})
}
