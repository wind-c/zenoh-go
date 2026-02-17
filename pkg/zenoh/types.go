// Package zenoh provides Go bindings for the zenoh protocol.
// It offers both owned and loaned types for interacting with zenoh.
//
// Memory Management Overview:
//
// Zenoh Go bindings follow the same ownership model as zenoh-c:
// - Owned types (OwnedXXX) own resources and MUST be explicitly dropped
// - Loaned types (XXX) are borrowed references and do not own resources
//
// Example usage:
//
//	config, err := zenoh.NewOwnedConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer config.Drop()  // Always defer drop
//
//	// Get a loaned reference from owned
//	loaned, err := zenoh.FromOwnedConfig(config)
//
// Key principles:
// 1. Always call Drop() on owned types when done
// 2. Use defer immediately after creating owned types
// 3. Check IsValid() before using any type
// 4. Loaned types become invalid when their owner is dropped
package zenoh

import (
	"errors"
	"log"
	"unsafe"
)

// =============================================================================
// Error Types
// =============================================================================

// ErrInvalidValue is returned when an operation is performed on an invalid owned type.
var ErrInvalidValue = errors.New("invalid zenoh value")

// ErrDropFailed is returned when dropping a resource fails.
var ErrDropFailed = errors.New("failed to drop zenoh resource")

// Result is the result of zenoh operations.
type Result struct {
	code int
}

// Error implements the error interface.
func (r Result) Error() string {
	if r.code == 0 {
		return ""
	}
	return errors.New("zenoh error").Error()
}

// IsOK returns true if the result is Z_OK.
func (r Result) IsOK() bool {
	return r.code == 0
}

// =============================================================================
// Owned Types
// =============================================================================

// OwnedConfig represents a zenoh configuration that owns its resources.
// It must be explicitly dropped to release underlying zenoh resources.
//
// # Memory Management
//
// OwnedConfig wraps a zenoh configuration handle. The configuration
// must be explicitly dropped to free the underlying zenoh resources.
//
// Note: The owned pointer (unsafe.Pointer) is used to store the C-owned
// config for passing to zenoh-c functions that require move semantics.
//
// Example:
//
//	config, err := zenoh.NewOwnedConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer config.Drop()
type OwnedConfig struct {
	// ptr is the underlying C pointer (stored as uintptr for type safety)
	ptr uintptr
	// owned is the owned config pointer for CGO move semantics
	owned unsafe.Pointer
}

// NewOwnedConfig creates a new owned zenoh configuration with default settings.
// The returned OwnedConfig must be dropped after use.
//
// Returns ErrInvalidValue if the configuration cannot be created.
func NewOwnedConfig() (*OwnedConfig, error) {
	// This function will be implemented in internal/cgo with actual zenoh-c calls
	// For now, return a placeholder that indicates the design intent
	return nil, errors.New("config creation requires zenoh-c bindings (see internal/cgo)")
}

// FromOwnedConfig creates a Config from an existing owned config.
// Note: This transfers ownership - the original OwnedConfig should not be used after this.
//
// The returned Config is a loaned reference that does not own resources.
func FromOwnedConfig(owned *OwnedConfig) (*Config, error) {
	if owned == nil || !owned.IsValid() {
		return nil, ErrInvalidValue
	}
	// This function will be implemented in internal/cgo
	return nil, errors.New("config loan requires zenoh-c bindings (see internal/cgo)")
}

// Drop releases the underlying zenoh configuration resources.
// After calling Drop, the OwnedConfig is invalidated and IsValid() returns false.
//
// It is safe to call Drop multiple times; subsequent calls are no-ops.
func (c *OwnedConfig) Drop() error {
	if c.ptr == 0 {
		return nil
	}
	// This function will be implemented in internal/cgo with actual z_drop call
	// For now, log the intent
	log.Print("[zenoh] config dropped")
	c.ptr = 0
	return nil
}

// IsValid returns true if the OwnedConfig contains a valid zenoh configuration.
// Returns false if the configuration is nil or has been dropped.
func (c *OwnedConfig) IsValid() bool {
	return c.ptr != 0
}

// Config is a loaned zenoh configuration reference.
// It borrows from an OwnedConfig and does not own resources.
//
// # Memory Management
//
// A Config is obtained from an OwnedConfig via FromOwnedConfig().
// The Config does not need to be dropped - it merely references
// the owned configuration. However, the Config becomes invalid
// if the OwnedConfig it was derived from is dropped.
type Config struct {
	// ptr is the underlying C pointer
	ptr uintptr
}

// IsValid returns true if the Config reference is valid.
// Returns false if the Config is nil or was derived from a dropped OwnedConfig.
func (c *Config) IsValid() bool {
	return c.ptr != 0
}

// =============================================================================
// Session Types
// =============================================================================

// SessionInfo contains information about a zenoh session.
type SessionInfo struct {
	ZID      []byte   // Zenoh ID (unique identifier)
	WhatAmI  string   // Role: "client", "router", or "peer"
	Locators []string // Connection addresses
}

// OwnedSession represents a zenoh session that owns its resources.
// It must be explicitly dropped to release underlying zenoh resources.
//
// # Memory Management
//
// The session is the main entry point for zenoh operations.
// It must be explicitly dropped to properly clean up all resources.
//
// Example:
//
//	session, err := zenoh.Open(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer session.Drop()
type OwnedSession struct {
	ptr   uintptr
	owned unsafe.Pointer
}

// NewOwnedSession creates a new owned zenoh session with default configuration.
// The returned OwnedSession must be dropped after use.
func NewOwnedSession() (*OwnedSession, error) {
	// Implemented in internal/cgo
	return nil, errors.New("session creation requires zenoh-c bindings (see internal/cgo)")
}

// FromOwnedSession creates a Session from an existing owned session.
// The returned Session is a loaned reference that does not own resources.
func FromOwnedSession(owned *OwnedSession) (*Session, error) {
	if owned == nil || !owned.IsValid() {
		return nil, ErrInvalidValue
	}
	// Implemented in internal/cgo
	return nil, errors.New("session loan requires zenoh-c bindings (see internal/cgo)")
}

// Drop releases the underlying zenoh session resources.
// After calling Drop, the OwnedSession is invalidated.
//
// It is safe to call Drop multiple times; subsequent calls are no-ops.
func (s *OwnedSession) Drop() error {
	if s.ptr == 0 {
		return nil
	}
	log.Print("[zenoh] session dropped")
	s.ptr = 0
	return nil
}

// IsValid returns true if the OwnedSession contains a valid zenoh session.
// Returns false if the session is nil or has been dropped.
func (s *OwnedSession) IsValid() bool {
	return s.ptr != 0
}

// Info returns session information including ZID, WhatAmI, and locators.
func (s *OwnedSession) Info() (*SessionInfo, error) {
	return nil, errors.New("Session.Info requires zenoh-c bindings")
}

// Session is a loaned zenoh session reference.
// It borrows from an OwnedSession and does not own resources.
//
// # Memory Management
//
// A Session is obtained from an OwnedSession via FromOwnedSession().
// The Session does not need to be dropped - it merely references
// the owned session. However, the Session becomes invalid
// if the OwnedSession it was derived from is dropped.
type Session struct {
	ptr uintptr
}

// IsValid returns true if the Session reference is valid.
// Returns false if the Session is nil or was derived from a dropped OwnedSession.
func (s *Session) IsValid() bool {
	return s.ptr != 0
}

// Close closes the zenoh session.
// Note: For OwnedSession, use Drop() instead.
func (s *Session) Close() error {
	if s.ptr == 0 {
		return nil
	}
	// Implemented in internal/cgo
	return errors.New("session close requires zenoh-c bindings (see internal/cgo)")
}

// =============================================================================
// KeyExpr Types
// =============================================================================

// OwnedKeyExpr represents a zenoh key expression that owns its resources.
// It must be explicitly dropped to release underlying zenoh resources.
//
// # Memory Management
//
// Key expressions are used to identify resources in zenoh.
// They must be dropped when no longer needed.
//
// Example:
//
//	keyExpr, err := zenoh.NewOwnedKeyExpr("demo/example/*")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer keyExpr.Drop()
type OwnedKeyExpr struct {
	ptr   uintptr
	owned unsafe.Pointer
}

// NewOwnedKeyExpr creates a new owned key expression from a string.
// The key expression string must be valid zenoh key expression syntax.
func NewOwnedKeyExpr(expr string) (*OwnedKeyExpr, error) {
	if expr == "" {
		return nil, errors.New("empty key expression")
	}
	// Implemented in internal/cgo
	return nil, errors.New("keyexpr creation requires zenoh-c bindings (see internal/cgo)")
}

// FromOwnedKeyExpr creates a KeyExpr from an existing owned key expression.
// The returned KeyExpr is a loaned reference that does not own resources.
func FromOwnedKeyExpr(owned *OwnedKeyExpr) (*KeyExpr, error) {
	if owned == nil || !owned.IsValid() {
		return nil, ErrInvalidValue
	}
	// Implemented in internal/cgo
	return nil, errors.New("keyexpr loan requires zenoh-c bindings (see internal/cgo)")
}

// Drop releases the underlying key expression resources.
// After calling Drop, the OwnedKeyExpr is invalidated.
func (k *OwnedKeyExpr) Drop() error {
	if k.ptr == 0 {
		return nil
	}
	log.Print("[zenoh] keyexpr dropped")
	k.ptr = 0
	return nil
}

// IsValid returns true if the OwnedKeyExpr is valid.
func (k *OwnedKeyExpr) IsValid() bool {
	return k.ptr != 0
}

// KeyExpr is a loaned key expression reference.
// It borrows from an OwnedKeyExpr and does not own resources.
type KeyExpr struct {
	ptr  uintptr
	expr string
}

// IsValid returns true if the KeyExpr reference is valid.
func (k *KeyExpr) IsValid() bool {
	if k == nil {
		return false
	}
	return k.ptr != 0 || k.expr != ""
}

// String returns the key expression as a string.
func (k *KeyExpr) String() string {
	if k == nil {
		return ""
	}
	if k.expr != "" {
		return k.expr
	}
	if k.ptr == 0 {
		return ""
	}
	return ""
}

// Expr returns the raw key expression string.
func (k *KeyExpr) Expr() string {
	if k == nil {
		return ""
	}
	return k.String()
}

// =============================================================================
// Publisher Types
// =============================================================================

// OwnedPublisher represents a zenoh publisher that owns its resources.
// It must be explicitly dropped to undeclare the publisher.
//
// # Memory Management
//
// A publisher is used to write data to a key expression.
// It must be dropped to undeclare the publisher and release resources.
//
// Example:
//
//	pub, err := zenoh.DeclarePublisher(session, keyExpr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pub.Drop()
type OwnedPublisher struct {
	ptr uintptr
}

// Drop releases the publisher by undeclareing it.
// After calling Drop, the OwnedPublisher is invalidated.
func (p *OwnedPublisher) Drop() error {
	if p.ptr == 0 {
		return nil
	}
	log.Print("[zenoh] publisher dropped")
	p.ptr = 0
	return nil
}

// IsValid returns true if the OwnedPublisher is valid.
func (p *OwnedPublisher) IsValid() bool {
	return p.ptr != 0
}

// Publisher is a loaned publisher reference.
// It borrows from an OwnedPublisher and does not own resources.
type Publisher struct {
	ptr uintptr
}

// IsValid returns true if the Publisher reference is valid.
func (p *Publisher) IsValid() bool {
	return p.ptr != 0
}

// =============================================================================
// Subscriber Types
// =============================================================================

// OwnedSubscriber represents a zenoh subscriber that owns its resources.
// It must be explicitly dropped to undeclare the subscriber.
//
// # Memory Management
//
// A subscriber is used to receive data from a key expression.
// It must be dropped to undeclare the subscriber and release resources.
//
// Example:
//
//	sub, err := zenoh.DeclareSubscriber(session, keyExpr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer sub.Drop()
type OwnedSubscriber struct {
	ptr uintptr
}

// Drop releases the subscriber by undeclareing it.
// After calling Drop, the OwnedSubscriber is invalidated.
func (s *OwnedSubscriber) Drop() error {
	if s == nil || s.ptr == 0 {
		return nil
	}
	log.Print("[zenoh] subscriber dropped")
	s.ptr = 0
	return nil
}

// IsValid returns true if the OwnedSubscriber is valid.
func (s *OwnedSubscriber) IsValid() bool {
	return s != nil && s.ptr != 0
}

// Undeclare is an alias for Drop.
func (s *OwnedSubscriber) Undeclare() error {
	return s.Drop()
}

// Subscriber is a loaned subscriber reference.
// It borrows from an OwnedSubscriber and does not own resources.
type Subscriber struct {
	ptr uintptr
}

// IsValid returns true if the Subscriber reference is valid.
func (s *Subscriber) IsValid() bool {
	return s != nil && s.ptr != 0
}

// =============================================================================
// Queryable Types
// =============================================================================

// OwnedQueryable represents a zenoh queryable that owns its resources.
// It must be explicitly dropped to undeclare the queryable.
//
// # Memory Management
//
// A queryable is used to respond to queries on a key expression.
// It must be dropped to undeclare the queryable and release resources.
//
// Example:
//
//	queryable, err := zenoh.DeclareQueryable(session, keyExpr, callback)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer queryable.Drop()
type OwnedQueryable struct {
	ptr    uintptr
	handle uintptr
}

// IsValid returns true if the OwnedQueryable is valid.
func (q *OwnedQueryable) IsValid() bool {
	return q != nil && q.ptr != 0
}

// Queryable is a loaned queryable reference.
// It borrows from an OwnedQueryable and does not own resources.
type Queryable struct {
	ptr uintptr
}

// =============================================================================
// Utility Functions
// =============================================================================

// Init initializes the zenoh runtime.
// This should be called before any other zenoh operations.
func Init() {
	log.Print("[zenoh] initializing zenoh runtime")
	// Zenoh runtime initialization is typically automatic
	// but we log the initialization for debugging purposes
}

// Log logs a message using the standard library log.
func Log(msg string) {
	log.Print("[zenoh] ", msg)
}

// =============================================================================
// Query Types
// =============================================================================

// QueryTarget specifies the target of a query.
type QueryTarget int

const (
	// QueryTargetBestMatching requests the best matching reply.
	QueryTargetBestMatching QueryTarget = iota
	// QueryTargetAll requests all matching replies.
	QueryTargetAll
	// QueryTargetComplete requests all replies with complete knowledge.
	QueryTargetComplete
)

// Consolidation specifies the consolidation mode for query replies.
type Consolidation int

const (
	// ConsolidationAuto lets zenoh decide the consolidation mode.
	ConsolidationAuto Consolidation = iota
	// ConsolidationNone requests no consolidation.
	ConsolidationNone
	// ConsolidationLatest requests only the latest values.
	ConsolidationLatest
	// ConsolidationMonotonic requests monotonic consolidation.
	ConsolidationMonotonic
)

// GetOptions contains options for Get operations.
type GetOptions struct {
	// Target specifies the query target.
	Target QueryTarget
	// Consolidation specifies the consolidation mode.
	Consolidation Consolidation
	// Timeout specifies the query timeout.
	Timeout int64 // in nanoseconds
}

// Reliability defines the reliability mode for pub/sub.
type Reliability int

const (
	// ReliabilityBestEffort does not guarantee delivery.
	ReliabilityBestEffort Reliability = 0
	// ReliabilityReliable guarantees delivery.
	ReliabilityReliable Reliability = 1
)

// CongestionControl defines how to handle network congestion.
type CongestionControl int

const (
	// CongestionControlBlock blocks on congestion (don't drop).
	CongestionControlBlock CongestionControl = 0
	// CongestionControlDrop drops messages on congestion.
	CongestionControlDrop CongestionControl = 1
)
