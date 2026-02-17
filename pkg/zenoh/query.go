package zenoh

import (
	"errors"
	"strings"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

var ErrInvalidSelector = errors.New("invalid selector")

type Reply struct {
	keyExpr  string
	value    []byte
	encoding *Encoding
	isOk     bool
	errMsg   string
	senderID []byte
	ptr      uintptr
}

func (r *Reply) KeyExpr() string {
	if r == nil {
		return ""
	}
	return r.keyExpr
}

func (r *Reply) Value() []byte {
	if r == nil {
		return nil
	}
	return r.value
}

func (r *Reply) Encoding() *Encoding {
	if r == nil {
		return nil
	}
	return r.encoding
}

func (r *Reply) IsOk() bool {
	if r == nil {
		return false
	}
	return r.isOk
}

func (r *Reply) Error() string {
	if r == nil {
		return ""
	}
	return r.errMsg
}

func (r *Reply) SenderID() []byte {
	if r == nil {
		return nil
	}
	return r.senderID
}

func (r *Reply) String() string {
	if r == nil {
		return "<nil>"
	}
	if !r.isOk {
		return "Reply{error: " + r.errMsg + "}"
	}
	return "Reply{keyExpr: " + r.keyExpr + ", value: " + string(r.value) + "}"
}

type ReplyCallback func(reply Reply)

type replyClosure struct {
	callback ReplyCallback
	channel  *ReplyChannel
}

type ReplyChannel struct {
	ch     chan Reply
	closed bool
}

func NewReplyChannel(bufferSize int) *ReplyChannel {
	if bufferSize <= 0 {
		bufferSize = 1
	}
	return &ReplyChannel{
		ch: make(chan Reply, bufferSize),
	}
}

func (r *ReplyChannel) Chan() <-chan Reply {
	return r.ch
}

func (r *ReplyChannel) Send(reply Reply) bool {
	if r.closed {
		return false
	}
	select {
	case r.ch <- reply:
		return true
	default:
		select {
		case <-r.ch:
			r.ch <- reply
			return true
		default:
			return false
		}
	}
}

func (r *ReplyChannel) Close() {
	if !r.closed {
		r.closed = true
		close(r.ch)
	}
}

func (r *ReplyChannel) IsClosed() bool {
	return r.closed
}

type ReplyIterator struct {
	ch      <-chan Reply
	current Reply
	valid   bool
}

func NewReplyIterator(ch <-chan Reply) *ReplyIterator {
	return &ReplyIterator{
		ch:    ch,
		valid: false,
	}
}

func (it *ReplyIterator) Next() bool {
	if it.ch == nil {
		return false
	}
	reply, ok := <-it.ch
	if !ok {
		return false
	}
	it.current = reply
	it.valid = true
	return true
}

func (it *ReplyIterator) Reply() Reply {
	return it.current
}

func (it *ReplyIterator) Valid() bool {
	return it.valid
}

func (it *ReplyIterator) KeyExpr() string {
	return it.current.KeyExpr()
}

func (it *ReplyIterator) Value() []byte {
	return it.current.Value()
}

func (it *ReplyIterator) Encoding() *Encoding {
	return it.current.Encoding()
}

func Get(session *OwnedSession, selector string, handler ReplyCallback) error {
	if session == nil || !session.IsValid() {
		return ErrInvalidValue
	}
	if selector == "" {
		return ErrInvalidSelector
	}
	if handler == nil {
		return errors.New("handler cannot be nil")
	}

	keyExpr, _, err := parseSelector(selector)
	if err != nil {
		return err
	}

	s := cgo.SessionFromOwnedPtr(session.ptr, session.owned)
	cb := func(data cgo.QueryReplyData) {
		reply := Reply{
			keyExpr: data.KeyExpr,
			value:   data.Payload,
			isOk:    data.Ok,
			errMsg:  data.ErrMsg,
			ptr:     0,
		}
		if data.Ok {
			reply.encoding = EncodingApplicationOctetStream
		}
		handler(reply)
	}
	return s.Get(keyExpr, cb)
}

func GetWithChannel(session *OwnedSession, selector string) (*ReplyChannel, error) {
	if session == nil || !session.IsValid() {
		return nil, ErrInvalidValue
	}
	if selector == "" {
		return nil, ErrInvalidSelector
	}

	keyExpr, _, err := parseSelector(selector)
	if err != nil {
		return nil, err
	}

	ch := NewReplyChannel(16)
	s := cgo.SessionFromOwnedPtr(session.ptr, session.owned)
	cb := func(data cgo.QueryReplyData) {
		reply := Reply{
			keyExpr: data.KeyExpr,
			value:   data.Payload,
			isOk:    data.Ok,
			errMsg:  data.ErrMsg,
			ptr:     0,
		}
		if data.Ok {
			reply.encoding = EncodingApplicationOctetStream
		}
		ch.Send(reply)
	}
	err = s.Get(keyExpr, cb)
	if err != nil {
		ch.Close()
		return nil, err
	}
	return ch, nil
}

func GetWithIterator(session *OwnedSession, selector string) (*ReplyIterator, error) {
	ch, err := GetWithChannel(session, selector)
	if err != nil {
		return nil, err
	}
	return NewReplyIterator(ch.Chan()), nil
}

func parseSelector(selector string) (keyExpr, params string, err error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return "", "", ErrInvalidSelector
	}

	idx := strings.Index(selector, "?")
	if idx == -1 {
		return selector, "", nil
	}

	if idx == 0 {
		return "", "", ErrInvalidSelector
	}

	keyExpr = selector[:idx]
	params = selector[idx+1:]
	return keyExpr, params, nil
}
