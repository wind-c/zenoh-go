package zenoh

import (
	"unsafe"

	"github.com/wind-c/zenoh-go/internal/cgo"
)

type OwnedShmProvider struct {
	ptr uintptr
}

func NewOwnedShmProvider(providerType, providerParams string) (*OwnedShmProvider, error) {
	return nil, ErrInvalidValue
}

func (s *OwnedSession) SharedMemoryProvider() (*ShmProvider, error) {
	if s == nil || !s.IsValid() {
		return nil, ErrInvalidValue
	}
	session := cgo.SessionFromOwnedPtr(s.ptr, s.owned)
	provider, err := session.ObtainShmProvider()
	if err != nil {
		return nil, err
	}
	return &ShmProvider{ptr: provider.Ptr}, nil
}

func (p *OwnedShmProvider) Drop() error {
	if p == nil || p.ptr == 0 {
		return nil
	}
	p.ptr = 0
	return nil
}

func (p *OwnedShmProvider) IsValid() bool {
	return p.ptr != 0
}

func (p *OwnedShmProvider) Check() bool {
	return false
}

func (p *OwnedShmProvider) Alloc(size int) (*OwnedShmBuf, error) {
	if p == nil || p.ptr == 0 {
		return nil, ErrInvalidValue
	}
	return nil, ErrInvalidValue
}

type ShmProvider struct {
	ptr uintptr
}

func (p *ShmProvider) IsValid() bool {
	return p.ptr != 0
}

type POSIXShmProvider struct {
	ptr unsafe.Pointer
}

func NewPOSIXShmProvider(layout string) (*POSIXShmProvider, error) {
	provider, err := cgo.NewPosixShmProvider(layout)
	if err != nil {
		return nil, err
	}
	return &POSIXShmProvider{ptr: unsafe.Pointer(provider.Ptr)}, nil
}

func (p *POSIXShmProvider) Alloc(size int) (*ShmBuf, error) {
	return nil, ErrInvalidValue
}

func (p *POSIXShmProvider) Drop() error {
	return nil
}

type OwnedShmBuf struct {
	ptr  uintptr
	data []byte
}

func NewOwnedShmBuf(data []byte) (*OwnedShmBuf, error) {
	return nil, ErrInvalidValue
}

func (b *OwnedShmBuf) Drop() error {
	if b == nil || b.ptr == 0 {
		return nil
	}
	b.ptr = 0
	b.data = nil
	return nil
}

func (b *OwnedShmBuf) IsValid() bool {
	return b.ptr != 0
}

func (b *OwnedShmBuf) Data() ([]byte, error) {
	if b == nil || b.ptr == 0 {
		return nil, ErrInvalidValue
	}
	return nil, ErrInvalidValue
}

func (b *OwnedShmBuf) Len() (int, error) {
	if b == nil || b.ptr == 0 {
		return 0, ErrInvalidValue
	}
	return 0, ErrInvalidValue
}

type ShmBuf struct {
	ptr  uintptr
	data []byte
}

func (b *ShmBuf) IsValid() bool {
	return b.ptr != 0
}
