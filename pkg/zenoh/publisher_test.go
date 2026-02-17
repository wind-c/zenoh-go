package zenoh

import (
	"errors"
	"testing"
)

func TestDeclarePublisherWithKeyExpr(t *testing.T) {
	tests := []struct {
		name    string
		session *OwnedSession
		keyExpr string
		wantErr bool
	}{
		{"nil session", nil, "demo/test", true},
		{"invalid session", &OwnedSession{ptr: 0}, "demo/test", true},
		{"empty keyExpr", &OwnedSession{ptr: 0}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeclarePublisherWithKeyExpr(tt.session, tt.keyExpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeclarePublisherWithKeyExpr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOwnedPublisher_Put(t *testing.T) {
	tests := []struct {
		name      string
		publisher *OwnedPublisher
		data      []byte
		encoding  *Encoding
		wantErr   bool
	}{
		{"nil publisher", nil, []byte("test"), EncodingTextPlain, true},
		{"invalid publisher", &OwnedPublisher{ptr: 0}, []byte("test"), EncodingTextPlain, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.publisher.Put(tt.data, tt.encoding)
			if (err != nil) != tt.wantErr {
				t.Errorf("OwnedPublisher.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOwnedPublisher_Delete(t *testing.T) {
	tests := []struct {
		name      string
		publisher *OwnedPublisher
		wantErr   bool
	}{
		{"nil publisher", nil, true},
		{"invalid publisher", &OwnedPublisher{ptr: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.publisher.Delete()
			if (err != nil) != tt.wantErr {
				t.Errorf("OwnedPublisher.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOwnedPublisher_Undeclare(t *testing.T) {
	t.Run("nil publisher", func(t *testing.T) {
		var p *OwnedPublisher
		err := p.Undeclare()
		if err != nil {
			t.Errorf("Undeclare() on nil publisher error = %v, want nil", err)
		}
	})

	t.Run("invalid publisher", func(t *testing.T) {
		p := &OwnedPublisher{ptr: 0}
		err := p.Undeclare()
		if err != nil {
			t.Errorf("Undeclare() on invalid publisher error = %v, want nil", err)
		}
	})
}

func TestOwnedPublisher_MatchingStatus(t *testing.T) {
	t.Run("nil publisher", func(t *testing.T) {
		var p *OwnedPublisher
		_, err := p.MatchingStatus()
		if !errors.Is(err, ErrInvalidPublisher) {
			t.Errorf("MatchingStatus() on nil publisher error = %v, want ErrInvalidPublisher", err)
		}
	})

	t.Run("invalid publisher", func(t *testing.T) {
		p := &OwnedPublisher{ptr: 0}
		_, err := p.MatchingStatus()
		if !errors.Is(err, ErrInvalidPublisher) {
			t.Errorf("MatchingStatus() on invalid publisher error = %v, want ErrInvalidPublisher", err)
		}
	})
}

func TestFromOwnedPublisher(t *testing.T) {
	tests := []struct {
		name    string
		owned   *OwnedPublisher
		wantErr bool
	}{
		{"nil owned", nil, true},
		{"invalid owned", &OwnedPublisher{ptr: 0}, true},
		// {"valid owned", &OwnedPublisher{ptr: 1}, false}, // Requires real zenoh publisher, would crash
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromOwnedPublisher(tt.owned)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromOwnedPublisher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanedPublisher_Put(t *testing.T) {
	tests := []struct {
		name      string
		publisher *Publisher
		data      []byte
		encoding  *Encoding
		wantErr   bool
	}{
		{"nil publisher", nil, []byte("test"), EncodingTextPlain, true},
		{"invalid publisher", &Publisher{ptr: 0}, []byte("test"), EncodingTextPlain, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.publisher.Put(tt.data, tt.encoding)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publisher.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanedPublisher_Delete(t *testing.T) {
	tests := []struct {
		name      string
		publisher *Publisher
		wantErr   bool
	}{
		{"nil publisher", nil, true},
		{"invalid publisher", &Publisher{ptr: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.publisher.Delete()
			if (err != nil) != tt.wantErr {
				t.Errorf("Publisher.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanedPublisher_MatchingStatus(t *testing.T) {
	t.Run("nil publisher", func(t *testing.T) {
		var p *Publisher
		_, err := p.MatchingStatus()
		if !errors.Is(err, ErrInvalidPublisher) {
			t.Errorf("MatchingStatus() on nil publisher error = %v, want ErrInvalidPublisher", err)
		}
	})

	t.Run("invalid publisher", func(t *testing.T) {
		p := &Publisher{ptr: 0}
		_, err := p.MatchingStatus()
		if !errors.Is(err, ErrInvalidPublisher) {
			t.Errorf("MatchingStatus() on invalid publisher error = %v, want ErrInvalidPublisher", err)
		}
	})
}

func TestMatchingStatus(t *testing.T) {
	m := &MatchingStatus{
		Matched: true,
	}

	if m.Matched != true {
		t.Errorf("Expected Matched=true, got %v", m.Matched)
	}
}

func TestPublisherOptions_Default(t *testing.T) {
	opts := DefaultPublisherOptions()
	if opts == nil {
		t.Fatal("DefaultPublisherOptions() returned nil")
	}
	if opts.Reliability != ReliabilityBestEffort {
		t.Errorf("Expected Reliability=ReliabilityBestEffort, got %v", opts.Reliability)
	}
	if opts.CongestionControl != CongestionControlDrop {
		t.Errorf("Expected CongestionControl=CongestionControlDrop, got %v", opts.CongestionControl)
	}
}

func TestPublisherOptions_Custom(t *testing.T) {
	opts := &PublisherOptions{
		Reliability:       ReliabilityReliable,
		CongestionControl: CongestionControlBlock,
	}
	if opts.Reliability != ReliabilityReliable {
		t.Errorf("Expected Reliability=ReliabilityReliable, got %v", opts.Reliability)
	}
	if opts.CongestionControl != CongestionControlBlock {
		t.Errorf("Expected CongestionControl=CongestionControlBlock, got %v", opts.CongestionControl)
	}
}
