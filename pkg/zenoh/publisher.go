package zenoh

import (
	"errors"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

var ErrInvalidPublisher = errors.New("invalid publisher")

type MatchingStatus struct {
	Matched bool
}

func DeclarePublisher(session *OwnedSession, keyExpr *OwnedKeyExpr) (*OwnedPublisher, error) {
	if session == nil || !session.IsValid() {
		return nil, ErrInvalidValue
	}
	if keyExpr == nil || !keyExpr.IsValid() {
		return nil, ErrInvalidValue
	}
	return nil, errors.New("DeclarePublisher requires OwnedKeyExpr with string representation")
}

func DeclarePublisherWithKeyExpr(session *OwnedSession, keyExpr string) (*OwnedPublisher, error) {
	if session == nil || !session.IsValid() {
		return nil, ErrInvalidValue
	}
	if keyExpr == "" {
		return nil, ErrInvalidKeyExpr
	}
	s := cgo.SessionFromOwnedPtr(session.ptr, session.owned)
	p, err := s.DeclarePublisherByKeyExpr(keyExpr)
	if err != nil {
		return nil, err
	}
	return &OwnedPublisher{ptr: p.Ptr}, nil
}

// PublisherOptions contains options for Publisher declaration.
type PublisherOptions struct {
	// Reliability specifies the reliability mode.
	Reliability Reliability
	// CongestionControl specifies the congestion control mode.
	CongestionControl CongestionControl
}

// DefaultPublisherOptions returns default publisher options.
func DefaultPublisherOptions() *PublisherOptions {
	return &PublisherOptions{
		Reliability:       ReliabilityBestEffort,
		CongestionControl: CongestionControlDrop,
	}
}

// DeclarePublisherWithOptions declares a publisher with custom options.
func DeclarePublisherWithOptions(session *OwnedSession, keyExpr string, opts *PublisherOptions) (*OwnedPublisher, error) {
	if session == nil || !session.IsValid() {
		return nil, ErrInvalidValue
	}
	if keyExpr == "" {
		return nil, ErrInvalidKeyExpr
	}
	if opts == nil {
		opts = DefaultPublisherOptions()
	}
	s := cgo.SessionFromOwnedPtr(session.ptr, session.owned)
	p, err := s.DeclarePublisherByKeyExprWithOptions(keyExpr, int(opts.Reliability), int(opts.CongestionControl))
	if err != nil {
		return nil, err
	}
	return &OwnedPublisher{ptr: p.Ptr}, nil
}

func (p *OwnedPublisher) Put(data []byte, encoding *Encoding) error {
	if p == nil || p.ptr == 0 {
		return ErrInvalidPublisher
	}
	pub := cgo.PublisherFromPtr(p.ptr)
	enc := encoding.toCGO()
	return pub.Put(data, enc)
}

func (p *OwnedPublisher) Delete() error {
	if p == nil || p.ptr == 0 {
		return ErrInvalidPublisher
	}
	pub := cgo.Publisher{Ptr: p.ptr}
	return pub.Delete()
}

func (p *OwnedPublisher) Undeclare() error {
	if p == nil || p.ptr == 0 {
		return nil
	}
	pub := cgo.Publisher{Ptr: p.ptr}
	err := pub.Undeclare()
	if err != nil {
		return err
	}
	p.ptr = 0
	return nil
}

func (p *OwnedPublisher) MatchingStatus() (*MatchingStatus, error) {
	if p == nil || p.ptr == 0 {
		return nil, ErrInvalidPublisher
	}
	pub := cgo.Publisher{Ptr: p.ptr}
	matched, err := pub.MatchingStatus()
	if err != nil {
		return nil, err
	}
	return &MatchingStatus{
		Matched: matched,
	}, nil
}

func FromOwnedPublisher(owned *OwnedPublisher) (*Publisher, error) {
	if owned == nil || !owned.IsValid() {
		return nil, ErrInvalidPublisher
	}
	return &Publisher{ptr: owned.ptr}, nil
}

func (p *Publisher) Put(data []byte, encoding *Encoding) error {
	if p == nil || p.ptr == 0 {
		return ErrInvalidPublisher
	}
	pub := cgo.Publisher{Ptr: p.ptr}
	enc := encoding.toCGO()
	return pub.Put(data, enc)
}

func (p *Publisher) Delete() error {
	if p == nil || p.ptr == 0 {
		return ErrInvalidPublisher
	}
	pub := cgo.Publisher{Ptr: p.ptr}
	return pub.Delete()
}

func (p *Publisher) MatchingStatus() (*MatchingStatus, error) {
	if p == nil || p.ptr == 0 {
		return nil, ErrInvalidPublisher
	}
	pub := cgo.Publisher{Ptr: p.ptr}
	matched, err := pub.MatchingStatus()
	if err != nil {
		return nil, err
	}
	return &MatchingStatus{
		Matched: matched,
	}, nil
}
