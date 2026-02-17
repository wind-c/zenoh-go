package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/include

#cgo windows LDFLAGS: -L${SRCDIR}/../../lib/windows -lzenohc -lws2_32 -liphlpapi -ladvapi32
#cgo linux LDFLAGS: -L${SRCDIR}/../../lib/linux -lzenohc
#cgo darwin LDFLAGS: -L${SRCDIR}/../../lib/darwin -lzenohc

#include <stdlib.h>
#include <string.h>
#include "zenoh.h"

extern void goSubscriberCallback(void *sample, void *context);

static void cSubscriberCallback(struct z_loaned_sample_t *sample, void *context) {
    goSubscriberCallback((void*)sample, context);
}

static void createClosureSample(struct z_owned_closure_sample_t *closure, void *context) {
    z_closure_sample(closure, cSubscriberCallback, NULL, context);
}

static void keyexprToString(const struct z_loaned_keyexpr_t *keyexpr, char *buf, size_t buf_len) {
    struct z_view_string_t view_str;
    z_keyexpr_as_view_string(keyexpr, &view_str);
    const struct z_loaned_string_t *loaned_str = z_view_string_loan(&view_str);
    size_t len = z_string_len(loaned_str);
    const char *data = z_string_data(loaned_str);
    if (len > 0 && len < buf_len) {
        memcpy(buf, data, len);
        buf[len] = '\0';
    } else if (buf_len > 0) {
        buf[0] = '\0';
    }
}

static size_t samplePayloadToSlice(const struct z_loaned_sample_t *sample, uint8_t *buf, size_t buf_len) {
    const struct z_loaned_bytes_t *payload = z_sample_payload(sample);
    if (payload == NULL) {
        return 0;
    }
    struct z_owned_slice_t slice;
    z_bytes_to_slice(payload, &slice);
    const struct z_loaned_slice_t *loaned = z_slice_loan(&slice);
    size_t len = z_slice_len(loaned);
    if (len > buf_len) {
        len = buf_len;
    }
    if (len > 0) {
        const uint8_t *data = z_slice_data(loaned);
        memcpy(buf, data, len);
    }
    z_slice_drop((struct z_moved_slice_t *)&slice);
    return len;
}

// Query Reply callback
extern void goReplyCallback(void *reply, void *context);

static void cReplyCallback(struct z_loaned_reply_t *reply, void *context) {
    goReplyCallback((void*)reply, context);
}

static void createClosureReply(struct z_owned_closure_reply_t *closure, void *context) {
    z_closure_reply(closure, cReplyCallback, NULL, context);
}

// Query callback
extern void goQueryCallback(void *query, void *context);

static void cQueryCallback(struct z_loaned_query_t *query, void *context) {
    goQueryCallback((void*)query, context);
}

static void createClosureQuery(struct z_owned_closure_query_t *closure, void *context) {
    z_closure_query(closure, cQueryCallback, NULL, context);
}

// Query helper functions
static void queryKeyexprToString(const struct z_loaned_query_t *query, char *buf, size_t buf_len) {
    const struct z_loaned_keyexpr_t *keyexpr = z_query_keyexpr(query);
    if (keyexpr == NULL || buf_len == 0) {
        if (buf_len > 0) buf[0] = '\0';
        return;
    }
    struct z_view_string_t view_str;
    z_keyexpr_as_view_string(keyexpr, &view_str);
    const struct z_loaned_string_t *loaned_str = z_view_string_loan(&view_str);
    size_t len = z_string_len(loaned_str);
    const char *data = z_string_data(loaned_str);
    if (len > 0 && len < buf_len) {
        memcpy(buf, data, len);
        buf[len] = '\0';
    } else if (buf_len > 0) {
        buf[0] = '\0';
    }
}

static size_t queryPayloadToSlice(const struct z_loaned_query_t *query, uint8_t *buf, size_t buf_len) {
    struct z_loaned_bytes_t *payload = z_query_payload(query);
    if (payload == NULL) {
        return 0;
    }
    struct z_owned_slice_t slice;
    z_bytes_to_slice(payload, &slice);
    const struct z_loaned_slice_t *loaned = z_slice_loan(&slice);
    size_t len = z_slice_len(loaned);
    if (len > buf_len) {
        len = buf_len;
    }
    if (len > 0) {
        const uint8_t *data = z_slice_data(loaned);
        memcpy(buf, data, len);
    }
    z_slice_drop((struct z_moved_slice_t *)&slice);
    return len;
}

// Query parameters helper
static size_t queryParametersToSlice(const struct z_loaned_query_t *query, char *buf, size_t buf_len) {
    struct z_view_string_t params;
    z_query_parameters(query, &params);
    const struct z_loaned_string_t *loaned_str = z_view_string_loan(&params);
    if (loaned_str == NULL || buf_len == 0) {
        if (buf_len > 0) buf[0] = '\0';
        return 0;
    }
    size_t len = z_string_len(loaned_str);
    const char *data = z_string_data(loaned_str);
    if (len > 0 && len < buf_len) {
        memcpy(buf, data, len);
        buf[len] = '\0';
    } else if (buf_len > 0) {
        buf[0] = '\0';
    }
    return len;
}

// Reply OK payload helper
static size_t replyOkPayloadToSlice(const struct z_loaned_reply_t *reply, uint8_t *buf, size_t buf_len) {
    struct z_loaned_sample_t *sample = z_reply_ok(reply);
    if (sample == NULL) {
        return 0;
    }
    const struct z_loaned_bytes_t *payload = z_sample_payload(sample);
    if (payload == NULL) {
        return 0;
    }
    struct z_owned_slice_t slice;
    z_bytes_to_slice(payload, &slice);
    const struct z_loaned_slice_t *loaned = z_slice_loan(&slice);
    size_t len = z_slice_len(loaned);
    if (len > buf_len) {
        len = buf_len;
    }
    if (len > 0) {
        const uint8_t *data = z_slice_data(loaned);
        memcpy(buf, data, len);
    }
    z_slice_drop((struct z_moved_slice_t *)&slice);
    return len;
}

// Reply OK keyexpr helper
static void replyOkKeyexprToString(const struct z_loaned_reply_t *reply, char *buf, size_t buf_len) {
    struct z_loaned_sample_t *sample = z_reply_ok(reply);
    if (sample == NULL || buf_len == 0) {
        if (buf_len > 0) buf[0] = '\0';
        return;
    }
    const struct z_loaned_keyexpr_t *keyexpr = z_sample_keyexpr(sample);
    if (keyexpr == NULL) {
        if (buf_len > 0) buf[0] = '\0';
        return;
    }
    struct z_view_string_t view_str;
    z_keyexpr_as_view_string(keyexpr, &view_str);
    const struct z_loaned_string_t *loaned_str = z_view_string_loan(&view_str);
    size_t len = z_string_len(loaned_str);
    const char *data = z_string_data(loaned_str);
    if (len > 0 && len < buf_len) {
        memcpy(buf, data, len);
        buf[len] = '\0';
    } else if (buf_len > 0) {
        buf[0] = '\0';
    }
}
*/
import "C"

import (
	"errors"
	"sync"
	"sync/atomic"
	"unsafe"
)

func Check(ret C.z_result_t) error {
	if ret != 0 {
		return errors.New("zenoh error")
	}
	return nil
}

type CallbackRegistry struct {
	mu       sync.RWMutex
	handlers map[uintptr]interface{}
	nextID   uintptr
}

func NewCallbackRegistry() *CallbackRegistry {
	return &CallbackRegistry{
		handlers: make(map[uintptr]interface{}),
	}
}

func (r *CallbackRegistry) Register(cb interface{}) uintptr {
	id := atomic.AddUintptr(&r.nextID, 1)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[id] = cb
	return id
}

func (r *CallbackRegistry) Unregister(id uintptr) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, id)
}

func (r *CallbackRegistry) Get(id uintptr) (interface{}, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cb, ok := r.handlers[id]
	return cb, ok
}

// Config
type Config struct {
	ptr   *C.z_loaned_config_t
	owned *C.z_owned_config_t
	Ptr   uintptr
}

// SetOwnedPtr sets the owned config pointer for move semantics
func (cfg *Config) SetOwnedPtr(ptr unsafe.Pointer) {
	if ptr != nil {
		cfg.owned = (*C.z_owned_config_t)(ptr)
	}
}

func ConfigDefault() (*Config, error) {
	var owned C.z_owned_config_t
	ret := C.z_config_default(&owned)
	if ret != 0 {
		return nil, Check(ret)
	}
	loaned := C.z_config_loan(&owned)
	return &Config{ptr: loaned, owned: &owned, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

func ConfigFromFile(path string) (*Config, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var owned C.z_owned_config_t
	err := C.zc_config_from_file(&owned, cPath)
	if err != 0 {
		return nil, Check(err)
	}
	loaned := C.z_config_loan(&owned)
	return &Config{ptr: loaned, owned: &owned, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

func ConfigFromStr(s string) (*Config, error) {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))

	var owned C.z_owned_config_t
	err := C.zc_config_from_str(&owned, cStr)
	if err != 0 {
		return nil, Check(err)
	}
	loaned := C.z_config_loan(&owned)
	return &Config{ptr: loaned, owned: &owned, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

func (cfg *Config) InsertJSON5(key, value string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	return Check(C.zc_config_insert_json5(cfg.ptr, cKey, cValue))
}

func (cfg *Config) OwnedPtr() unsafe.Pointer {
	if cfg.owned == nil {
		return nil
	}
	return unsafe.Pointer(cfg.owned)
}

func (cfg *Config) SetLoaned(ptr unsafe.Pointer) {
	if cfg.ptr == nil && ptr != nil {
		cfg.ptr = (*C.z_loaned_config_t)(ptr)
	}
}

func (cfg *Config) SetOwned(ptr unsafe.Pointer) {
	if cfg.owned == nil && ptr != nil {
		cfg.owned = (*C.z_owned_config_t)(ptr)
	}
}

func (cfg *Config) Open() (*Session, error) {
	var ownedSession C.z_owned_session_t
	if cfg.owned == nil {
		return nil, errors.New("config not initialized properly")
	}
	err := C.z_open(&ownedSession, (*C.z_moved_config_t)(unsafe.Pointer(cfg.owned)), nil)
	if err != 0 {
		return nil, Check(err)
	}
	loaned := C.z_session_loan(&ownedSession)
	return &Session{ptr: loaned, owned: &ownedSession, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

func (cfg *Config) Drop() {
	if cfg.ptr != nil {
		cfg.ptr = nil
	}
	if cfg.owned != nil {
		C.z_config_drop((*C.z_moved_config_t)(unsafe.Pointer(cfg.owned)))
		cfg.owned = nil
	}
}

// Session
type Session struct {
	ptr   *C.z_loaned_session_t
	owned *C.z_owned_session_t
	Ptr   uintptr
}

func SessionFromPtr(ptr uintptr) *Session {
	return &Session{
		ptr: (*C.z_loaned_session_t)(unsafe.Pointer(ptr)),
		Ptr: ptr,
	}
}

func SessionFromOwnedPtr(ptr uintptr, owned unsafe.Pointer) *Session {
	return &Session{
		ptr:   (*C.z_loaned_session_t)(unsafe.Pointer(ptr)),
		owned: (*C.z_owned_session_t)(owned),
		Ptr:   ptr,
	}
}

func (s *Session) OwnedPtr() unsafe.Pointer {
	return unsafe.Pointer(s.owned)
}

func (s *Session) Close() error {
	if s.owned != nil {
		C.z_session_drop((*C.z_moved_session_t)(unsafe.Pointer(s.owned)))
		s.owned = nil
		s.ptr = nil
	}
	return nil
}

func (s *Session) DeclareKeyExpr(keyExpr string) (*OwnedKeyExpr, error) {
	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var owned C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&owned, cKeyExpr)
	if err != 0 {
		return nil, Check(err)
	}
	loaned := C.z_keyexpr_loan(&owned)
	return &OwnedKeyExpr{owned: owned, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

func (s *Session) DeclarePublisherByKeyExpr(keyExpr string) (*Publisher, error) {
	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var ownedKeyExpr C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&ownedKeyExpr, cKeyExpr)
	if err != 0 {
		return nil, Check(err)
	}
	loanedKeyExpr := C.z_keyexpr_loan(&ownedKeyExpr)

	var owned C.z_owned_publisher_t
	err = C.z_declare_publisher(s.ptr, &owned, loanedKeyExpr, nil)
	C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&ownedKeyExpr)))
	if err != 0 {
		return nil, Check(err)
	}
	loaned := C.z_publisher_loan(&owned)
	return &Publisher{ptr: loaned, owned: &owned, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

func (s *Session) DeclarePublisherByKeyExprWithOptions(keyExpr string, reliability int, congestion int) (*Publisher, error) {
	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var ownedKeyExpr C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&ownedKeyExpr, cKeyExpr)
	if err != 0 {
		return nil, Check(err)
	}
	loanedKeyExpr := C.z_keyexpr_loan(&ownedKeyExpr)

	var opts C.z_publisher_options_t
	C.z_publisher_options_default(&opts)
	opts.reliability = C.enum_z_reliability_t(reliability)
	opts.congestion_control = C.enum_z_congestion_control_t(congestion)

	var owned C.z_owned_publisher_t
	ret := C.z_declare_publisher(s.ptr, &owned, loanedKeyExpr, &opts)
	C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&ownedKeyExpr)))
	if ret != 0 {
		return nil, Check(ret)
	}
	loaned := C.z_publisher_loan(&owned)
	return &Publisher{ptr: loaned, owned: &owned, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

// KeyExpr
type OwnedKeyExpr struct {
	owned C.z_owned_keyexpr_t
	Ptr   uintptr
}

func OwnedKeyExprFromPtr(ptr uintptr) *OwnedKeyExpr {
	return &OwnedKeyExpr{
		Ptr: ptr,
	}
}

func OwnedKeyExprFromOwnedPtr(ptr uintptr, owned unsafe.Pointer) *OwnedKeyExpr {
	return &OwnedKeyExpr{
		owned: *(*C.z_owned_keyexpr_t)(owned),
		Ptr:   ptr,
	}
}

func (k *OwnedKeyExpr) OwnedPtr() unsafe.Pointer {
	return unsafe.Pointer(&k.owned)
}

func (k *OwnedKeyExpr) Drop() {
	C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&k.owned)))
}

// Publisher
type Publisher struct {
	ptr   *C.z_loaned_publisher_t
	owned *C.z_owned_publisher_t
	Ptr   uintptr
}

func PublisherFromPtr(ptr uintptr) *Publisher {
	return &Publisher{
		ptr: (*C.z_loaned_publisher_t)(unsafe.Pointer(ptr)),
		Ptr: ptr,
	}
}

func (p *Publisher) Put(payload []byte, encoding *Encoding) error {
	if len(payload) == 0 {
		return nil
	}

	var ownedBytes C.z_owned_bytes_t
	cPayload := C.CBytes(payload)
	defer C.free(cPayload)

	ret := C.z_bytes_copy_from_buf(&ownedBytes, (*C.uint8_t)(cPayload), C.size_t(len(payload)))
	if ret != 0 {
		return Check(ret)
	}

	var opts C.z_publisher_put_options_t
	C.z_publisher_put_options_default(&opts)

	moveResult := C.z_publisher_put(p.ptr, (*C.z_moved_bytes_t)(unsafe.Pointer(&ownedBytes)), &opts)
	return Check(moveResult)
}

func (p *Publisher) Delete() error {
	var opts C.z_publisher_delete_options_t
	C.z_publisher_delete_options_default(&opts)
	return Check(C.z_publisher_delete(p.ptr, &opts))
}

func (p *Publisher) Undeclare() error {
	if p.owned != nil {
		ret := C.z_undeclare_publisher((*C.z_moved_publisher_t)(unsafe.Pointer(p.owned)))
		if ret != 0 {
			return Check(ret)
		}
		p.owned = nil
		p.ptr = nil
	}
	return nil
}

func (p *Publisher) MatchingStatus() (bool, error) {
	var status C.z_matching_status_t
	C.z_publisher_get_matching_status(p.ptr, &status)
	return bool(status.matching), nil
}

// Subscriber
type SubscriberCallback func(SampleData)

type SampleData struct {
	KeyExpr string
	Payload []byte
}

var subscriberRegistry = NewCallbackRegistry()

//export goSubscriberCallback
func goSubscriberCallback(sample unsafe.Pointer, context unsafe.Pointer) {
	handle := uintptr(context)
	cb, ok := subscriberRegistry.Get(handle)
	if !ok {
		return
	}
	callback, ok := cb.(SubscriberCallback)
	if !ok {
		return
	}

	var keyExprBuf [256]byte
	C.keyexprToString(C.z_sample_keyexpr((*C.z_loaned_sample_t)(sample)), (*C.char)(unsafe.Pointer(&keyExprBuf)), C.size_t(256))
	keyExpr := C.GoStringN((*C.char)(unsafe.Pointer(&keyExprBuf)), C.int(C.strlen((*C.char)(unsafe.Pointer(&keyExprBuf)))))

	var payloadBuf [4096]byte
	payloadLen := C.samplePayloadToSlice((*C.z_loaned_sample_t)(sample), (*C.uint8_t)(unsafe.Pointer(&payloadBuf)), C.size_t(4096))
	payload := C.GoBytes(unsafe.Pointer(&payloadBuf), C.int(payloadLen))

	callback(SampleData{
		KeyExpr: keyExpr,
		Payload: payload,
	})
}

func (s *Session) DeclareSubscriber(keyExpr string, callback SubscriberCallback) (*Subscriber, error) {
	if callback == nil {
		return nil, errors.New("callback cannot be nil")
	}

	handle := subscriberRegistry.Register(callback)

	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var ownedKeyExpr C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&ownedKeyExpr, cKeyExpr)
	if err != 0 {
		subscriberRegistry.Unregister(handle)
		return nil, Check(err)
	}
	defer C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&ownedKeyExpr)))

	loanedKeyExpr := C.z_keyexpr_loan(&ownedKeyExpr)

	var closure C.z_owned_closure_sample_t
	C.createClosureSample(&closure, unsafe.Pointer(handle))

	var opts C.z_subscriber_options_t
	C.z_subscriber_options_default(&opts)

	var ownedSubscriber C.z_owned_subscriber_t
	ret := C.z_declare_subscriber(s.ptr, &ownedSubscriber, loanedKeyExpr, (*C.z_moved_closure_sample_t)(unsafe.Pointer(&closure)), &opts)
	if ret != 0 {
		subscriberRegistry.Unregister(handle)
		return nil, Check(ret)
	}

	loaned := C.z_subscriber_loan(&ownedSubscriber)
	return &Subscriber{ptr: loaned, owned: &ownedSubscriber, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

func (s *Session) DeclareSubscriberWithOptions(keyExpr string, callback SubscriberCallback, reliability int) (*Subscriber, error) {
	if callback == nil {
		return nil, errors.New("callback cannot be nil")
	}

	handle := subscriberRegistry.Register(callback)

	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var ownedKeyExpr C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&ownedKeyExpr, cKeyExpr)
	if err != 0 {
		subscriberRegistry.Unregister(handle)
		return nil, Check(err)
	}
	defer C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&ownedKeyExpr)))

	loanedKeyExpr := C.z_keyexpr_loan(&ownedKeyExpr)

	var closure C.z_owned_closure_sample_t
	C.createClosureSample(&closure, unsafe.Pointer(handle))

	var opts C.z_subscriber_options_t
	C.z_subscriber_options_default(&opts)

	var ownedSubscriber C.z_owned_subscriber_t
	ret := C.z_declare_subscriber(s.ptr, &ownedSubscriber, loanedKeyExpr, (*C.z_moved_closure_sample_t)(unsafe.Pointer(&closure)), &opts)
	if ret != 0 {
		subscriberRegistry.Unregister(handle)
		return nil, Check(ret)
	}

	loaned := C.z_subscriber_loan(&ownedSubscriber)
	return &Subscriber{ptr: loaned, owned: &ownedSubscriber, Ptr: uintptr(unsafe.Pointer(loaned))}, nil
}

type Subscriber struct {
	ptr   *C.z_loaned_subscriber_t
	owned *C.z_owned_subscriber_t
	Ptr   uintptr
}

func (s *Subscriber) Undeclare() error {
	if s.owned != nil {
		ret := C.z_undeclare_subscriber((*C.z_moved_subscriber_t)(unsafe.Pointer(s.owned)))
		if ret != 0 {
			return Check(ret)
		}
		s.owned = nil
		s.ptr = nil
	}
	return nil
}

// Encoding (simplified - just a placeholder)
type Encoding struct {
	ptr *C.z_loaned_encoding_t
}

// Query Reply types
type QueryReplyCallback func(QueryReplyData)

type QueryReplyData struct {
	Ok      bool
	KeyExpr string
	Payload []byte
	ErrMsg  string
}

var replyRegistry = NewCallbackRegistry()

//export goReplyCallback
func goReplyCallback(reply unsafe.Pointer, context unsafe.Pointer) {
	handle := uintptr(context)
	cb, ok := replyRegistry.Get(handle)
	if !ok {
		return
	}
	callback, ok := cb.(QueryReplyCallback)
	if !ok {
		return
	}

	isOk := C.z_reply_is_ok((*C.z_loaned_reply_t)(reply))

	var data QueryReplyData
	if isOk {
		var keyExprBuf [256]byte
		C.replyOkKeyexprToString((*C.z_loaned_reply_t)(reply), (*C.char)(unsafe.Pointer(&keyExprBuf)), C.size_t(256))
		keyExpr := C.GoStringN((*C.char)(unsafe.Pointer(&keyExprBuf)), C.int(C.strlen((*C.char)(unsafe.Pointer(&keyExprBuf)))))

		var payloadBuf [4096]byte
		payloadLen := C.replyOkPayloadToSlice((*C.z_loaned_reply_t)(reply), (*C.uint8_t)(unsafe.Pointer(&payloadBuf)), C.size_t(4096))
		payload := C.GoBytes(unsafe.Pointer(&payloadBuf), C.int(payloadLen))

		data = QueryReplyData{
			Ok:      true,
			KeyExpr: keyExpr,
			Payload: payload,
		}
	} else {
		replyErr := C.z_reply_err((*C.z_loaned_reply_t)(reply))
		if replyErr != nil {
			errPayload := C.z_reply_err_payload(replyErr)
			if errPayload != nil {
				var errSlice C.z_owned_slice_t
				C.z_bytes_to_slice(errPayload, &errSlice)
				errLoaned := C.z_slice_loan(&errSlice)
				errLen := C.z_slice_len(errLoaned)
				if errLen > 0 {
					errData := C.z_slice_data(errLoaned)
					data = QueryReplyData{
						Ok:     false,
						ErrMsg: C.GoStringN((*C.char)(unsafe.Pointer(errData)), C.int(errLen)),
					}
				} else {
					data = QueryReplyData{
						Ok:     false,
						ErrMsg: "unknown error",
					}
				}
				C.z_slice_drop((*C.z_moved_slice_t)(unsafe.Pointer(&errSlice)))
			} else {
				data = QueryReplyData{
					Ok:     false,
					ErrMsg: "unknown error",
				}
			}
		} else {
			data = QueryReplyData{
				Ok:     false,
				ErrMsg: "unknown error",
			}
		}
	}

	C.z_reply_drop((*C.z_moved_reply_t)(unsafe.Pointer(reply)))

	callback(data)
}

func (s *Session) Get(keyExpr string, callback QueryReplyCallback) error {
	if callback == nil {
		return errors.New("callback cannot be nil")
	}

	handle := replyRegistry.Register(callback)

	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var ownedKeyExpr C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&ownedKeyExpr, cKeyExpr)
	if err != 0 {
		replyRegistry.Unregister(handle)
		return Check(err)
	}
	defer C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&ownedKeyExpr)))

	loanedKeyExpr := C.z_keyexpr_loan(&ownedKeyExpr)

	var closure C.z_owned_closure_reply_t
	C.createClosureReply(&closure, unsafe.Pointer(handle))

	var opts C.z_get_options_t
	C.z_get_options_default(&opts)

	ret := C.z_get(s.ptr, loanedKeyExpr, nil, (*C.z_moved_closure_reply_t)(unsafe.Pointer(&closure)), &opts)
	if ret != 0 {
		replyRegistry.Unregister(handle)
		return Check(ret)
	}

	return nil
}

// Query types
type Query struct {
	ptr        *C.z_loaned_query_t
	KeyExpr    string
	Parameters string
	Payload    []byte
}

func (q *Query) Reply(keyExpr string, payload []byte) error {
	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var ownedKeyExpr C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&ownedKeyExpr, cKeyExpr)
	if err != 0 {
		return Check(err)
	}
	loanedKeyExpr := C.z_keyexpr_loan(&ownedKeyExpr)
	defer C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&ownedKeyExpr)))

	var ownedBytes C.z_owned_bytes_t
	if len(payload) > 0 {
		cPayload := C.CBytes(payload)
		defer C.free(cPayload)
		ret := C.z_bytes_copy_from_buf(&ownedBytes, (*C.uint8_t)(cPayload), C.size_t(len(payload)))
		if ret != 0 {
			return Check(ret)
		}
	}

	var opts C.z_query_reply_options_t
	C.z_query_reply_options_default(&opts)

	return Check(C.z_query_reply(q.ptr, loanedKeyExpr, (*C.z_moved_bytes_t)(unsafe.Pointer(&ownedBytes)), &opts))
}

func (q *Query) ReplyErr(errMsg string) error {
	var ownedBytes C.z_owned_bytes_t
	if len(errMsg) > 0 {
		cMsg := C.CBytes([]byte(errMsg))
		defer C.free(cMsg)
		ret := C.z_bytes_copy_from_buf(&ownedBytes, (*C.uint8_t)(cMsg), C.size_t(len(errMsg)))
		if ret != 0 {
			return Check(ret)
		}
	}

	var opts C.z_query_reply_err_options_t
	C.z_query_reply_err_options_default(&opts)

	return Check(C.z_query_reply_err(q.ptr, (*C.z_moved_bytes_t)(unsafe.Pointer(&ownedBytes)), &opts))
}

type QueryableCallback func(Query)

//export goQueryCallback
func goQueryCallback(query unsafe.Pointer, context unsafe.Pointer) {
	handle := uintptr(context)
	cb, ok := queryableRegistry.Get(handle)
	if !ok {
		return
	}
	callback, ok := cb.(QueryableCallback)
	if !ok {
		return
	}

	var keyExprBuf [256]byte
	C.queryKeyexprToString((*C.z_loaned_query_t)(query), (*C.char)(unsafe.Pointer(&keyExprBuf)), C.size_t(256))
	keyExpr := C.GoStringN((*C.char)(unsafe.Pointer(&keyExprBuf)), C.int(C.strlen((*C.char)(unsafe.Pointer(&keyExprBuf)))))

	var paramsBuf [1024]byte
	C.queryParametersToSlice((*C.z_loaned_query_t)(query), (*C.char)(unsafe.Pointer(&paramsBuf)), C.size_t(1024))
	params := C.GoStringN((*C.char)(unsafe.Pointer(&paramsBuf)), C.int(C.strlen((*C.char)(unsafe.Pointer(&paramsBuf)))))

	var payloadBuf [4096]byte
	payloadLen := C.queryPayloadToSlice((*C.z_loaned_query_t)(query), (*C.uint8_t)(unsafe.Pointer(&payloadBuf)), C.size_t(4096))
	payload := C.GoBytes(unsafe.Pointer(&payloadBuf), C.int(payloadLen))

	callback(Query{
		ptr:        (*C.z_loaned_query_t)(query),
		KeyExpr:    keyExpr,
		Parameters: params,
		Payload:    payload,
	})
}

var queryableRegistry = NewCallbackRegistry()

func (s *Session) DeclareQueryable(keyExpr string, callback QueryableCallback) (*Queryable, error) {
	if callback == nil {
		return nil, errors.New("callback cannot be nil")
	}

	handle := queryableRegistry.Register(callback)

	cKeyExpr := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(cKeyExpr))

	var ownedKeyExpr C.z_owned_keyexpr_t
	err := C.z_keyexpr_from_str(&ownedKeyExpr, cKeyExpr)
	if err != 0 {
		queryableRegistry.Unregister(handle)
		return nil, Check(err)
	}
	defer C.z_keyexpr_drop((*C.z_moved_keyexpr_t)(unsafe.Pointer(&ownedKeyExpr)))

	loanedKeyExpr := C.z_keyexpr_loan(&ownedKeyExpr)

	var closure C.z_owned_closure_query_t
	C.createClosureQuery(&closure, unsafe.Pointer(handle))

	var opts C.z_queryable_options_t
	C.z_queryable_options_default(&opts)

	var ownedQueryable C.z_owned_queryable_t
	ret := C.z_declare_queryable(s.ptr, &ownedQueryable, loanedKeyExpr, (*C.z_moved_closure_query_t)(unsafe.Pointer(&closure)), &opts)
	if ret != 0 {
		queryableRegistry.Unregister(handle)
		return nil, Check(ret)
	}

	loaned := C.z_queryable_loan(&ownedQueryable)
	return &Queryable{session: s, ptr: loaned, owned: &ownedQueryable, Ptr: uintptr(unsafe.Pointer(loaned)), handle: handle}, nil
}

type Queryable struct {
	session *Session
	ptr     *C.z_loaned_queryable_t
	owned   *C.z_owned_queryable_t
	Ptr     uintptr
	handle  uintptr
}

func QueryableFromPtr(ptr uintptr) *Queryable {
	return &Queryable{
		ptr: (*C.z_loaned_queryable_t)(unsafe.Pointer(ptr)),
		Ptr: ptr,
	}
}

func (q *Queryable) Handle() uintptr {
	return q.handle
}

func (q *Queryable) Undeclare() error {
	if q.owned != nil {
		C.z_queryable_drop((*C.z_moved_queryable_t)(unsafe.Pointer(q.owned)))
		queryableRegistry.Unregister(q.handle)
		q.owned = nil
		q.ptr = nil
	}
	return nil
}

// =============================================================================
// Shared Memory
// =============================================================================

type ShmProvider struct {
	ptr   *C.z_loaned_shm_provider_t
	owned *C.z_owned_shm_provider_t
	Ptr   uintptr
}

func (s *Session) ObtainShmProvider() (*ShmProvider, error) {
	return nil, errors.New("shared memory not available in this zenoh-c build")
}

type PosixShmProvider struct {
	ptr   *C.z_loaned_shm_provider_t
	owned *C.z_owned_shm_provider_t
	Ptr   uintptr
}

func NewPosixShmProvider(layout string) (*PosixShmProvider, error) {
	return nil, errors.New("shared memory not available in this zenoh-c build")
}

func (p *PosixShmProvider) Alloc(size int) (*ShmBuf, error) {
	return nil, errors.New("shared memory not available in this zenoh-c build")
}

func (p *PosixShmProvider) Drop() error {
	return nil
}

type ShmBuf struct {
	ptr *C.z_loaned_shm_t
	Ptr uintptr
}

func (b *ShmBuf) Data() []byte {
	return nil
}

func (b *ShmBuf) Len() int {
	return 0
}

func (b *ShmBuf) Drop() {
}
