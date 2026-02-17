package zenoh

import (
	"errors"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

var ErrInvalidQueryable = errors.New("invalid queryable")
var ErrInvalidQuery = errors.New("invalid query")

type Query struct {
	keyExpr    string
	parameters string
	payload    []byte
	ptr        uintptr
	cgoQuery   *cgo.Query
}

func (q *Query) KeyExpr() string {
	if q == nil {
		return ""
	}
	return q.keyExpr
}

func (q *Query) Parameters() string {
	if q == nil {
		return ""
	}
	return q.parameters
}

func (q *Query) Value() []byte {
	if q == nil {
		return nil
	}
	return q.payload
}

func (q *Query) Reply(keyExpr string, payload []byte, encoding *Encoding) error {
	if q == nil || (q.ptr == 0 && q.cgoQuery == nil) {
		return ErrInvalidQuery
	}
	if q.cgoQuery != nil {
		return q.cgoQuery.Reply(keyExpr, payload)
	}
	return errors.New("Query.Reply requires cgo query")
}

func (q *Query) ReplyErr(payload []byte) error {
	if q == nil || (q.ptr == 0 && q.cgoQuery == nil) {
		return ErrInvalidQuery
	}
	if q.cgoQuery != nil {
		errMsg := string(payload)
		return q.cgoQuery.ReplyErr(errMsg)
	}
	return errors.New("Query.ReplyErr requires cgo query")
}

type QueryCallback func(query Query)

type queryableClosure struct {
	callback QueryCallback
	channel  *QueryChannel
}

type QueryChannel struct {
	ch     chan Query
	closed bool
}

func NewQueryChannel(bufferSize int) *QueryChannel {
	if bufferSize <= 0 {
		bufferSize = 1
	}
	return &QueryChannel{
		ch: make(chan Query, bufferSize),
	}
}

func (r *QueryChannel) Chan() <-chan Query {
	return r.ch
}

func (r *QueryChannel) Send(query Query) bool {
	if r.closed {
		return false
	}
	select {
	case r.ch <- query:
		return true
	default:
		select {
		case <-r.ch:
			r.ch <- query
			return true
		default:
			return false
		}
	}
}

func (r *QueryChannel) Close() {
	if !r.closed {
		r.closed = true
		close(r.ch)
	}
}

func (r *QueryChannel) IsClosed() bool {
	return r.closed
}

func (q *Queryable) IsValid() bool {
	return q != nil && q.ptr != 0
}

func (q *OwnedQueryable) Undeclare() error {
	if q == nil || q.ptr == 0 {
		return nil
	}
	qable := cgo.QueryableFromPtr(q.ptr)
	err := qable.Undeclare()
	if err != nil {
		return err
	}
	q.ptr = 0
	q.handle = 0
	return nil
}

func (q *OwnedQueryable) Drop() error {
	return q.Undeclare()
}

func DeclareQueryable(session *OwnedSession, keyExpr string, callback QueryCallback) (*OwnedQueryable, error) {
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
	cgoCallback := func(cgoQuery cgo.Query) {
		query := Query{
			keyExpr:    cgoQuery.KeyExpr,
			parameters: cgoQuery.Parameters,
			payload:    cgoQuery.Payload,
			cgoQuery:   &cgoQuery,
		}
		callback(query)
	}
	qable, err := s.DeclareQueryable(keyExpr, cgoCallback)
	if err != nil {
		return nil, err
	}
	return &OwnedQueryable{ptr: qable.Ptr, handle: qable.Handle()}, nil
}
