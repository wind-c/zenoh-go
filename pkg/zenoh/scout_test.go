package zenoh

import (
	"testing"
	"time"
)

func TestWhatAmI_String(t *testing.T) {
	tests := []struct {
		name     string
		whatAmI  WhatAmI
		expected string
	}{
		{"router", WhatAmIRouter, "router"},
		{"peer", WhatAmIPeer, "peer"},
		{"client", WhatAmIClient, "client"},
		{"unknown", WhatAmI(255), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.whatAmI.String(); got != tt.expected {
				t.Errorf("WhatAmI.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultScoutOptions(t *testing.T) {
	opts := DefaultScoutOptions()
	if opts == nil {
		t.Fatal("DefaultScoutOptions() returned nil")
	}
	if opts.WhatAmI != WhatAmIPeer {
		t.Errorf("WhatAmI = %v, want %v", opts.WhatAmI, WhatAmIPeer)
	}
	if opts.Timeout != 0 {
		t.Errorf("Timeout = %v, want 0", opts.Timeout)
	}
	if opts.StartAsync != false {
		t.Errorf("StartAsync = %v, want false", opts.StartAsync)
	}
}

func TestHello_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		hello   *Hello
		wantErr bool
	}{
		{"nil hello", nil, false},
		{"valid hello", &Hello{ZID: []byte("test"), WhatAmI: WhatAmIPeer, Locators: []string{"tcp/127.0.0.1:7447"}}, true},
		{"empty zid", &Hello{ZID: []byte{}, WhatAmI: WhatAmIPeer}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hello.IsValid(); got != tt.wantErr {
				t.Errorf("Hello.IsValid() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestHello_String(t *testing.T) {
	tests := []struct {
		name    string
		hello   *Hello
		want    string
		wantErr bool
	}{
		{"nil hello", nil, "Hello(nil)", false},
		{"valid hello", &Hello{ZID: []byte("test"), WhatAmI: WhatAmIPeer, Locators: []string{"tcp/127.0.0.1:7447"}}, "Hello{zid:test, whatami:peer, locators:[tcp/127.0.0.1:7447]}", false},
		{"empty locators", &Hello{ZID: []byte("test"), WhatAmI: WhatAmIRouter, Locators: []string{}}, "Hello{zid:test, whatami:router, locators:[]}", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hello.String()
			if got != tt.want {
				t.Errorf("Hello.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScout(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		_, err := Scout(WhatAmIPeer, nil)
		if err == nil {
			t.Error("Scout() expected error (no zenoh-c bindings)")
		}
	})

	t.Run("with options", func(t *testing.T) {
		opts := DefaultScoutOptions()
		opts.Timeout = time.Second
		_, err := Scout(WhatAmIPeer, opts)
		if err == nil {
			t.Error("Scout() expected error (no zenoh-c bindings)")
		}
	})

	t.Run("router target", func(t *testing.T) {
		_, err := Scout(WhatAmIRouter, nil)
		if err == nil {
			t.Error("Scout() expected error (no zenoh-c bindings)")
		}
	})

	t.Run("client target", func(t *testing.T) {
		_, err := Scout(WhatAmIClient, nil)
		if err == nil {
			t.Error("Scout() expected error (no zenoh-c bindings)")
		}
	})
}

func TestScoutBlocking(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		_, err := ScoutBlocking(WhatAmIPeer, nil)
		if err == nil {
			t.Error("ScoutBlocking() expected error (no zenoh-c bindings)")
		}
	})

	t.Run("with timeout", func(t *testing.T) {
		opts := &ScoutOptions{
			WhatAmI: WhatAmIPeer,
			Timeout: 500 * time.Millisecond,
		}
		_, err := ScoutBlocking(WhatAmIPeer, opts)
		if err == nil {
			t.Error("ScoutBlocking() expected error (no zenoh-c bindings)")
		}
	})
}

func TestScoutOptions_Validation(t *testing.T) {
	t.Run("set whatami", func(t *testing.T) {
		opts := &ScoutOptions{
			WhatAmI:    WhatAmIRouter,
			Timeout:    5 * time.Second,
			StartAsync: true,
		}
		if opts.WhatAmI != WhatAmIRouter {
			t.Errorf("WhatAmI = %v, want %v", opts.WhatAmI, WhatAmIRouter)
		}
		if opts.Timeout != 5*time.Second {
			t.Errorf("Timeout = %v, want 5s", opts.Timeout)
		}
		if !opts.StartAsync {
			t.Error("StartAsync should be true")
		}
	})
}

func TestHello_EmptyZID(t *testing.T) {
	hello := &Hello{
		ZID:      []byte{},
		WhatAmI:  WhatAmIPeer,
		Locators: []string{},
	}
	if hello.IsValid() {
		t.Error("Hello with empty ZID should not be valid")
	}
}

func TestScout_InvalidWhatAmI(t *testing.T) {
	opts := DefaultScoutOptions()
	opts.WhatAmI = WhatAmIClient

	_, err := Scout(WhatAmIClient, opts)
	if err == nil {
		t.Error("Scout() expected error (no zenoh-c bindings)")
	}
}
