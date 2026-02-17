package zenoh

import (
	"errors"
	"time"
)

// =============================================================================
// WhatAmI - Scout Target Types
// =============================================================================

// WhatAmI represents the type of zenoh entity to scout for.
type WhatAmI uint8

const (
	// WhatAmIRouter scouts for routers.
	WhatAmIRouter WhatAmI = iota
	// WhatAmIPeer scouts for peers.
	WhatAmIPeer
	// WhatAmIClient scouts for clients.
	WhatAmIClient
)

// String returns the string representation of WhatAmI.
func (w WhatAmI) String() string {
	switch w {
	case WhatAmIRouter:
		return "router"
	case WhatAmIPeer:
		return "peer"
	case WhatAmIClient:
		return "client"
	default:
		return "unknown"
	}
}

// =============================================================================
// Scout Options
// =============================================================================

// ScoutOptions contains options for the Scout operation.
// This is equivalent to z_scout_options_t in zenoh-c.
type ScoutOptions struct {
	// WhatAmI defines the type of entities to scout for (router, peer, or client).
	WhatAmI WhatAmI

	// Timeout is the duration to scout for.
	// If 0, scouting runs indefinitely until cancelled.
	Timeout time.Duration

	// StartAsync enables asynchronous scouting (callback-based).
	// If false, scouting returns a channel of Hello messages.
	StartAsync bool
}

// DefaultScoutOptions returns default scouting options.
func DefaultScoutOptions() *ScoutOptions {
	return &ScoutOptions{
		WhatAmI:    WhatAmIPeer,
		Timeout:    0,
		StartAsync: false,
	}
}

// =============================================================================
// Hello - Discovery Result
// =============================================================================

// Hello represents the information received from a discovered peer or router.
// This is equivalent to z_hello_t in zenoh-c.
//
// A Hello message is received in response to scouting and contains
// information about a discovered zenoh entity.
type Hello struct {
	// ZID is the Zenoh ID of the discovered entity.
	ZID []byte

	// WhatAmI is the type of the discovered entity (router, peer, or client).
	WhatAmI WhatAmI

	// Locators are the connection locators of the discovered entity.
	Locators []string
}

// IsValid returns true if the Hello contains valid information.
func (h *Hello) IsValid() bool {
	return h != nil && len(h.ZID) > 0
}

// String returns a string representation of the Hello.
func (h *Hello) String() string {
	if h == nil {
		return "Hello(nil)"
	}
	return "Hello{zid:" + string(h.ZID) + ", whatami:" + h.WhatAmI.String() + ", locators:" + formatLocators(h.Locators) + "}"
}

func formatLocators(locators []string) string {
	if len(locators) == 0 {
		return "[]"
	}
	result := "["
	for i, loc := range locators {
		if i > 0 {
			result += ", "
		}
		result += loc
	}
	return result + "]"
}

// =============================================================================
// Scout - Dynamic Discovery
// =============================================================================

// Scout performs dynamic discovery of peers and/or routers in the zenoh network.
// This is equivalent to z_scout() in zenoh-c.
func Scout(whatami WhatAmI, opts *ScoutOptions) (<-chan *Hello, error) {
	return nil, errors.New("Scout requires zenoh-c bindings")
}

// ScoutBlocking performs blocking discovery of peers and/or routers.
func ScoutBlocking(whatami WhatAmI, opts *ScoutOptions) ([]*Hello, error) {
	return nil, errors.New("ScoutBlocking requires zenoh-c bindings")
}
