package zenoh

import (
	"testing"
)

func TestGet_Validation(t *testing.T) {
	tests := []struct {
		name     string
		session  *OwnedSession
		selector string
		handler  ReplyCallback
		wantErr  bool
	}{
		{"nil session", nil, "demo/test/*", func(r Reply) {}, true},
		{"invalid session", &OwnedSession{ptr: 0}, "demo/test/*", func(r Reply) {}, true},
		{"empty selector", &OwnedSession{ptr: 0}, "", func(r Reply) {}, true},
		{"nil handler", &OwnedSession{ptr: 0}, "demo/test/*", nil, true},
		{"valid params", &OwnedSession{ptr: 0}, "demo/test/*", func(r Reply) {}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Get(tt.session, tt.selector, tt.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetWithChannel_Validation(t *testing.T) {
	tests := []struct {
		name     string
		session  *OwnedSession
		selector string
		wantErr  bool
	}{
		{"nil session", nil, "demo/test/*", true},
		{"invalid session", &OwnedSession{ptr: 0}, "demo/test/*", true},
		{"empty selector", &OwnedSession{ptr: 0}, "", true},
		{"valid params", &OwnedSession{ptr: 0}, "demo/test/*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetWithChannel(tt.session, tt.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWithChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetWithIterator_Validation(t *testing.T) {
	tests := []struct {
		name     string
		session  *OwnedSession
		selector string
		wantErr  bool
	}{
		{"nil session", nil, "demo/test/*", true},
		{"invalid session", &OwnedSession{ptr: 0}, "demo/test/*", true},
		{"empty selector", &OwnedSession{ptr: 0}, "", true},
		{"valid params", &OwnedSession{ptr: 0}, "demo/test/*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetWithIterator(tt.session, tt.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWithIterator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReply(t *testing.T) {
	r := &Reply{
		keyExpr:  "demo/test/key1",
		value:    []byte("hello"),
		encoding: EncodingTextPlain,
		isOk:     true,
		errMsg:   "",
		senderID: []byte{1, 2, 3},
	}

	if r.KeyExpr() != "demo/test/key1" {
		t.Errorf("Expected KeyExpr=demo/test/key1, got %s", r.KeyExpr())
	}
	if string(r.Value()) != "hello" {
		t.Errorf("Expected Value=hello, got %s", string(r.Value()))
	}
	if r.Encoding() != EncodingTextPlain {
		t.Errorf("Expected Encoding=EncodingTextPlain, got %v", r.Encoding())
	}
	if !r.IsOk() {
		t.Error("IsOk() should return true")
	}
	if r.Error() != "" {
		t.Errorf("Expected Error='', got %s", r.Error())
	}

	str := r.String()
	if str == "" {
		t.Error("Reply.String() should not be empty")
	}
}

func TestReply_Error(t *testing.T) {
	r := &Reply{
		keyExpr:  "",
		value:    nil,
		encoding: nil,
		isOk:     false,
		errMsg:   "not found",
		senderID: nil,
	}

	if r.IsOk() {
		t.Error("IsOk() should return false for error reply")
	}
	if r.Error() != "not found" {
		t.Errorf("Expected Error='not found', got %s", r.Error())
	}

	str := r.String()
	if str == "" {
		t.Error("Reply.String() should not be empty")
	}
}

func TestReply_Nil(t *testing.T) {
	var r *Reply

	if r.KeyExpr() != "" {
		t.Error("Nil Reply KeyExpr should return empty string")
	}
	if r.Value() != nil {
		t.Error("Nil Reply Value should return nil")
	}
	if r.Encoding() != nil {
		t.Error("Nil Reply Encoding should return nil")
	}
	if r.IsOk() {
		t.Error("Nil Reply IsOk should return false")
	}
	if r.Error() != "" {
		t.Error("Nil Reply Error should return empty string")
	}
	if r.SenderID() != nil {
		t.Error("Nil Reply SenderID should return nil")
	}
}

func TestReplyChannel(t *testing.T) {
	t.Run("create with zero size", func(t *testing.T) {
		ch := NewReplyChannel(0)
		if ch == nil {
			t.Error("NewReplyChannel(0) should not return nil")
		}
		if ch.ch == nil {
			t.Error("channel should not be nil")
		}
	})

	t.Run("create with positive size", func(t *testing.T) {
		ch := NewReplyChannel(10)
		if cap(ch.ch) != 10 {
			t.Errorf("Expected capacity 10, got %d", cap(ch.ch))
		}
	})

	t.Run("send and receive", func(t *testing.T) {
		ch := NewReplyChannel(10)
		reply := Reply{keyExpr: "test", value: []byte("data"), isOk: true}

		if !ch.Send(reply) {
			t.Error("Send() should succeed")
		}

		select {
		case received := <-ch.ch:
			if received.KeyExpr() != "test" {
				t.Errorf("Expected KeyExpr=test, got %s", received.KeyExpr())
			}
		default:
			t.Error("Should have received reply")
		}
	})

	t.Run("close", func(t *testing.T) {
		ch := NewReplyChannel(10)
		ch.Close()
		if !ch.IsClosed() {
			t.Error("IsClosed() should return true after Close()")
		}
	})

	t.Run("send after close", func(t *testing.T) {
		ch := NewReplyChannel(10)
		ch.Close()
		reply := Reply{keyExpr: "test", value: []byte("data"), isOk: true}
		if ch.Send(reply) {
			t.Error("Send() should return false after Close()")
		}
	})
}

func TestReplyIterator(t *testing.T) {
	t.Run("nil channel", func(t *testing.T) {
		it := NewReplyIterator(nil)
		if it == nil {
			t.Error("NewReplyIterator(nil) should not return nil")
		}
		if it.Next() {
			t.Error("Next() on nil channel should return false")
		}
	})

	t.Run("empty channel", func(t *testing.T) {
		ch := make(chan Reply)
		it := NewReplyIterator(ch)
		go func() {
			close(ch)
		}()

		if it.Next() {
			t.Error("Next() on empty channel should return false")
		}
	})

	t.Run("iterate multiple replies", func(t *testing.T) {
		ch := make(chan Reply, 3)
		ch <- Reply{keyExpr: "key1", value: []byte("val1"), isOk: true}
		ch <- Reply{keyExpr: "key2", value: []byte("val2"), isOk: true}
		ch <- Reply{keyExpr: "key3", value: []byte("val3"), isOk: true}
		close(ch)

		it := NewReplyIterator(ch)
		count := 0
		for it.Next() {
			count++
			if !it.Valid() {
				t.Error("Valid() should return true after Next() returns true")
			}
			_ = it.KeyExpr()
			_ = it.Value()
			_ = it.Encoding()
		}

		if count != 3 {
			t.Errorf("Expected 3 replies, got %d", count)
		}
	})

	t.Run("error reply", func(t *testing.T) {
		ch := make(chan Reply, 1)
		ch <- Reply{keyExpr: "", value: nil, isOk: false, errMsg: "not found"}
		close(ch)

		it := NewReplyIterator(ch)
		if !it.Next() {
			t.Error("Expected at least one reply")
		}
		reply := it.Reply()
		if reply.IsOk() {
			t.Error("Expected error reply")
		}
		if reply.Error() != "not found" {
			t.Errorf("Expected error 'not found', got '%s'", reply.Error())
		}
	})
}
