package zenoh

import (
	"testing"
)

func TestQueryChannel(t *testing.T) {
	ch := NewQueryChannel(10)
	if ch == nil {
		t.Fatal("NewQueryChannel returned nil")
	}

	if ch.IsClosed() {
		t.Error("channel should not be closed initially")
	}

	ch.Send(Query{keyExpr: "test/key", payload: []byte("test value")})

	select {
	case q := <-ch.Chan():
		if q.keyExpr != "test/key" {
			t.Errorf("expected keyExpr 'test/key', got '%s'", q.keyExpr)
		}
		if string(q.payload) != "test value" {
			t.Errorf("expected payload 'test value', got '%s'", string(q.payload))
		}
	default:
		t.Error("expected to receive query from channel")
	}

	ch.Close()
	if !ch.IsClosed() {
		t.Error("channel should be closed after Close()")
	}
}

func TestQueryChannel_Closed(t *testing.T) {
	ch := NewQueryChannel(1)
	ch.Close()

	result := ch.Send(Query{keyExpr: "test"})
	if result {
		t.Error("Send should return false on closed channel")
	}
}

func TestQueryKeyExpr(t *testing.T) {
	q := Query{
		keyExpr: "demo/example/123",
		payload: []byte("hello"),
	}

	if q.KeyExpr() != "demo/example/123" {
		t.Errorf("expected keyExpr 'demo/example/123', got '%s'", q.KeyExpr())
	}
}

func TestQueryValue(t *testing.T) {
	q := Query{
		keyExpr: "demo/example/123",
		payload: []byte("hello"),
	}

	if string(q.Value()) != "hello" {
		t.Errorf("expected payload 'hello', got '%s'", string(q.Value()))
	}
}

func TestQueryNil(t *testing.T) {
	var q *Query

	if q.KeyExpr() != "" {
		t.Error("nil query should return empty keyExpr")
	}

	if q.Value() != nil {
		t.Error("nil query should return nil value")
	}
}

func TestOwnedQueryableNil(t *testing.T) {
	var q *OwnedQueryable

	if q.IsValid() {
		t.Error("nil OwnedQueryable should not be valid")
	}

	if q.Drop() != nil {
		t.Error("Drop on nil should return nil error")
	}

	if q.Undeclare() != nil {
		t.Error("Undeclare on nil should return nil error")
	}
}

func TestQueryableNil(t *testing.T) {
	var q *Queryable

	if q.IsValid() {
		t.Error("nil Queryable should not be valid")
	}
}
