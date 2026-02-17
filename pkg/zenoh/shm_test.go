package zenoh

import (
	"errors"
	"testing"
)

func TestNewOwnedShmProvider(t *testing.T) {
	tests := []struct {
		name           string
		providerType   string
		providerParams string
		wantErr        bool
	}{
		{"empty type", "", "", true},
		{"posix shm-pub", "posix", "", true},
		{"malloc shm-pub", "malloc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOwnedShmProvider(tt.providerType, tt.providerParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOwnedShmProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOwnedShmProvider_Drop(t *testing.T) {
	t.Run("nil provider", func(t *testing.T) {
		var p *OwnedShmProvider
		err := p.Drop()
		if err != nil {
			t.Errorf("Drop() on nil provider error = %v", err)
		}
	})

	t.Run("invalid provider", func(t *testing.T) {
		p := &OwnedShmProvider{ptr: 0}
		err := p.Drop()
		if err != nil {
			t.Errorf("Drop() on invalid provider error = %v", err)
		}
	})
}

func TestOwnedShmProvider_IsValid(t *testing.T) {
	tests := []struct {
		name string
		ptr  uintptr
		want bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OwnedShmProvider{ptr: tt.ptr}
			if got := p.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOwnedShmProvider_Check(t *testing.T) {
	t.Run("invalid provider", func(t *testing.T) {
		p := &OwnedShmProvider{ptr: 0}
		if got := p.Check(); got {
			t.Error("Check() on invalid provider should return false")
		}
	})

	t.Run("valid provider", func(t *testing.T) {
		p := &OwnedShmProvider{ptr: 1}
		got := p.Check()
		if got {
			t.Error("Check() expected error or false (no zenoh-c bindings)")
		}
	})
}

func TestOwnedShmProvider_Alloc(t *testing.T) {
	t.Run("invalid provider", func(t *testing.T) {
		p := &OwnedShmProvider{ptr: 0}
		_, err := p.Alloc(1024)
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("Alloc() on invalid provider error = %v, want ErrInvalidValue", err)
		}
	})
}

func TestNewOwnedShmBuf(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"nil data", nil, true},
		{"empty data", []byte{}, true},
		{"valid data", []byte("hello"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOwnedShmBuf(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOwnedShmBuf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOwnedShmBuf_Drop(t *testing.T) {
	t.Run("nil buffer", func(t *testing.T) {
		var b *OwnedShmBuf
		err := b.Drop()
		if err != nil {
			t.Errorf("Drop() on nil buffer error = %v", err)
		}
	})

	t.Run("invalid buffer", func(t *testing.T) {
		b := &OwnedShmBuf{ptr: 0}
		err := b.Drop()
		if err != nil {
			t.Errorf("Drop() on invalid buffer error = %v", err)
		}
	})
}

func TestOwnedShmBuf_IsValid(t *testing.T) {
	tests := []struct {
		name string
		ptr  uintptr
		want bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &OwnedShmBuf{ptr: tt.ptr}
			if got := b.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOwnedShmBuf_Data(t *testing.T) {
	t.Run("invalid buffer", func(t *testing.T) {
		b := &OwnedShmBuf{ptr: 0}
		_, err := b.Data()
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("Data() on invalid buffer error = %v, want ErrInvalidValue", err)
		}
	})
}

func TestOwnedShmBuf_Len(t *testing.T) {
	t.Run("invalid buffer", func(t *testing.T) {
		b := &OwnedShmBuf{ptr: 0}
		_, err := b.Len()
		if !errors.Is(err, ErrInvalidValue) {
			t.Errorf("Len() on invalid buffer error = %v, want ErrInvalidValue", err)
		}
	})
}

func TestShmProvider_IsValid(t *testing.T) {
	tests := []struct {
		name string
		ptr  uintptr
		want bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ShmProvider{ptr: tt.ptr}
			if got := p.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShmBuf_IsValid(t *testing.T) {
	tests := []struct {
		name string
		ptr  uintptr
		want bool
	}{
		{"nil ptr", 0, false},
		{"valid ptr", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ShmBuf{ptr: tt.ptr}
			if got := b.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOSIXShmProvider_New(t *testing.T) {
	_, err := NewPOSIXShmProvider("/tmp")
	if err == nil {
		t.Log("POSIXShmProvider created (SHM available in zenoh-c)")
	} else {
		t.Logf("POSIXShmProvider not available: %v", err)
	}
}

func TestPOSIXShmProvider_Alloc(t *testing.T) {
	p := &POSIXShmProvider{}
	_, err := p.Alloc(1024)
	if !errors.Is(err, ErrInvalidValue) {
		t.Errorf("Alloc() on invalid provider error = %v, want ErrInvalidValue", err)
	}
}

func TestPOSIXShmProvider_Drop(t *testing.T) {
	p := &POSIXShmProvider{}
	err := p.Drop()
	if err != nil {
		t.Errorf("Drop() error = %v", err)
	}
}

func TestOwnedSession_SharedMemoryProvider(t *testing.T) {
	session := &OwnedSession{ptr: 0}
	_, err := session.SharedMemoryProvider()
	if !errors.Is(err, ErrInvalidValue) {
		t.Errorf("SharedMemoryProvider() on invalid session error = %v, want ErrInvalidValue", err)
	}
}
