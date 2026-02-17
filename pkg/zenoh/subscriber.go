package zenoh

import (
	"errors"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

// ErrInvalidSubscriber is returned when an operation is performed on an invalid subscriber.
var ErrInvalidSubscriber = errors.New("invalid subscriber")

// Sample represents a zenoh sample received from a subscription.
type Sample struct {
	KeyExpr  string
	Payload  []byte
	Encoding *Encoding
}

func (s *Sample) String() string {
	if s == nil {
		return "<nil>"
	}
	return "Sample{keyExpr: " + s.KeyExpr + ", payload: " + string(s.Payload) + "}"
}

// RingChannel is a ring buffer channel for receiving subscriber samples.
type RingChannel struct {
	ch     chan Sample
	closed bool
}

func NewRingChannel(bufferSize int) *RingChannel {
	if bufferSize <= 0 {
		bufferSize = 1
	}
	return &RingChannel{
		ch: make(chan Sample, bufferSize),
	}
}

func (r *RingChannel) Chan() <-chan Sample {
	return r.ch
}

func (r *RingChannel) Send(sample Sample) bool {
	if r.closed {
		return false
	}
	select {
	case r.ch <- sample:
		return true
	default:
		select {
		case <-r.ch:
			r.ch <- sample
			return true
		default:
			return false
		}
	}
}

func (r *RingChannel) Close() {
	if !r.closed {
		r.closed = true
		close(r.ch)
	}
}

func (r *RingChannel) IsClosed() bool {
	return r.closed
}

// SubscriberCallback is a function type that handles received samples.
type SubscriberCallback func(sample Sample)

// DeclareSubscriber declares a subscriber for the given key expression with a callback.
func DeclareSubscriber(session *OwnedSession, keyExpr string, callback SubscriberCallback) (*OwnedSubscriber, error) {
	if session == nil || !session.IsValid() {
		return nil, ErrInvalidValue
	}
	if keyExpr == "" {
		return nil, ErrInvalidKeyExpr
	}
	if callback == nil {
		return nil, errors.New("callback cannot be nil")
	}

	s := cgo.SessionFromOwnedPtr(session.ptr, session.owned)

	cgoCallback := func(sample cgo.SampleData) {
		callback(Sample{
			KeyExpr:  sample.KeyExpr,
			Payload:  sample.Payload,
			Encoding: nil,
		})
	}

	sub, err := s.DeclareSubscriber(keyExpr, cgoCallback)
	if err != nil {
		return nil, err
	}

	return &OwnedSubscriber{ptr: sub.Ptr}, nil
}

// SubscriberOptions contains options for Subscriber declaration.
type SubscriberOptions struct {
	// Reliability specifies the reliability mode.
	Reliability Reliability
}

// DefaultSubscriberOptions returns default subscriber options.
func DefaultSubscriberOptions() *SubscriberOptions {
	return &SubscriberOptions{
		Reliability: ReliabilityReliable,
	}
}

// DeclareSubscriberWithOptions declares a subscriber with custom options.
func DeclareSubscriberWithOptions(session *OwnedSession, keyExpr string, callback SubscriberCallback, opts *SubscriberOptions) (*OwnedSubscriber, error) {
	if session == nil || !session.IsValid() {
		return nil, ErrInvalidValue
	}
	if keyExpr == "" {
		return nil, ErrInvalidKeyExpr
	}
	if callback == nil {
		return nil, errors.New("callback cannot be nil")
	}
	if opts == nil {
		opts = DefaultSubscriberOptions()
	}

	s := cgo.SessionFromOwnedPtr(session.ptr, session.owned)

	cgoCallback := func(sample cgo.SampleData) {
		callback(Sample{
			KeyExpr:  sample.KeyExpr,
			Payload:  sample.Payload,
			Encoding: nil,
		})
	}

	sub, err := s.DeclareSubscriberWithOptions(keyExpr, cgoCallback, int(opts.Reliability))
	if err != nil {
		return nil, err
	}

	return &OwnedSubscriber{ptr: sub.Ptr}, nil
}
