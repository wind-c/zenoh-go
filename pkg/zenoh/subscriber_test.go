package zenoh

import (
	"testing"
)

func TestDeclareSubscriberWithCallback(t *testing.T) {
	tests := []struct {
		name     string
		session  *OwnedSession
		keyExpr  string
		callback SubscriberCallback
		wantErr  bool
	}{
		{"nil session", nil, "demo/test", func(s Sample) {}, true},
		{"invalid session", &OwnedSession{ptr: 0}, "demo/test", func(s Sample) {}, true},
		{"empty keyExpr", &OwnedSession{ptr: 0}, "", func(s Sample) {}, true},
		{"nil callback", &OwnedSession{ptr: 0}, "demo/test", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeclareSubscriber(tt.session, tt.keyExpr, tt.callback)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeclareSubscriber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOwnedSubscriberUndeclare_Call(t *testing.T) {
	t.Run("nil subscriber", func(t *testing.T) {
		var s *OwnedSubscriber
		err := s.Undeclare()
		if err != nil {
			t.Errorf("Undeclare() on nil subscriber error = %v, want nil", err)
		}
	})

	t.Run("invalid subscriber", func(t *testing.T) {
		s := &OwnedSubscriber{ptr: 0}
		err := s.Undeclare()
		if err != nil {
			t.Errorf("Undeclare() on invalid subscriber error = %v, want nil", err)
		}
	})
}

func TestOwnedSubscriber_Validation(t *testing.T) {
	tests := []struct {
		name       string
		subscriber *OwnedSubscriber
		want       bool
	}{
		{"nil subscriber", nil, false},
		{"invalid subscriber", &OwnedSubscriber{ptr: 0}, false},
		{"valid subscriber", &OwnedSubscriber{ptr: 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subscriber.IsValid(); got != tt.want {
				t.Errorf("OwnedSubscriber.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscriber_Validation(t *testing.T) {
	tests := []struct {
		name       string
		subscriber *Subscriber
		want       bool
	}{
		{"nil subscriber", nil, false},
		{"invalid subscriber", &Subscriber{ptr: 0}, false},
		{"valid subscriber", &Subscriber{ptr: 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subscriber.IsValid(); got != tt.want {
				t.Errorf("Subscriber.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSample(t *testing.T) {
	s := &Sample{
		KeyExpr:  "demo/test",
		Payload:  []byte("hello"),
		Encoding: EncodingTextPlain,
	}

	if s.KeyExpr != "demo/test" {
		t.Errorf("Expected KeyExpr=demo/test, got %s", s.KeyExpr)
	}
	if string(s.Payload) != "hello" {
		t.Errorf("Expected Payload=hello, got %s", s.Payload)
	}
	if s.Encoding != EncodingTextPlain {
		t.Errorf("Expected Encoding=EncodingTextPlain, got %v", s.Encoding)
	}

	str := s.String()
	if str == "" {
		t.Error("Sample.String() should not be empty")
	}
}

func TestRingChannel(t *testing.T) {
	t.Run("create with zero size", func(t *testing.T) {
		ch := NewRingChannel(0)
		if ch == nil {
			t.Error("NewRingChannel(0) should not return nil")
		}
		if ch.ch == nil {
			t.Error("channel should not be nil")
		}
	})

	t.Run("create with positive size", func(t *testing.T) {
		ch := NewRingChannel(10)
		if cap(ch.ch) != 10 {
			t.Errorf("Expected capacity 10, got %d", cap(ch.ch))
		}
	})

	t.Run("send and receive", func(t *testing.T) {
		ch := NewRingChannel(10)
		sample := Sample{KeyExpr: "test", Payload: []byte("data")}

		if !ch.Send(sample) {
			t.Error("Send() should succeed")
		}

		select {
		case received := <-ch.ch:
			if received.KeyExpr != "test" {
				t.Errorf("Expected KeyExpr=test, got %s", received.KeyExpr)
			}
		default:
			t.Error("Should have received sample")
		}
	})

	t.Run("close", func(t *testing.T) {
		ch := NewRingChannel(10)
		ch.Close()
		if !ch.IsClosed() {
			t.Error("IsClosed() should return true after Close()")
		}
	})

	t.Run("send after close", func(t *testing.T) {
		ch := NewRingChannel(10)
		ch.Close()
		sample := Sample{KeyExpr: "test", Payload: []byte("data")}
		if ch.Send(sample) {
			t.Error("Send() should return false after Close()")
		}
	})
}

func TestSubscriberOptions_Default(t *testing.T) {
	opts := DefaultSubscriberOptions()
	if opts == nil {
		t.Fatal("DefaultSubscriberOptions() returned nil")
	}
	if opts.Reliability != ReliabilityReliable {
		t.Errorf("Expected Reliability=ReliabilityReliable, got %v", opts.Reliability)
	}
}

func TestSubscriberOptions_Custom(t *testing.T) {
	opts := &SubscriberOptions{
		Reliability: ReliabilityBestEffort,
	}
	if opts.Reliability != ReliabilityBestEffort {
		t.Errorf("Expected Reliability=ReliabilityBestEffort, got %v", opts.Reliability)
	}
}
