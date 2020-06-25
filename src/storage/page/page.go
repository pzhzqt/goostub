package page

import (
	"common"
	"unsafe"
)

const (
	sizePageHeader  = 8
	offsetPageStart = 0
	offsetLSN       = 4
)

type Page interface {
	GetData() []byte
	GetPageID() common.PageID
	GetPinCount() int
	IsDirty() bool
	GetLSN() common.LSN
	SetLSN(common.LSN)
	WLatch()
	WUnlatch()
	RLatch()
	RUnlatch()
}

// Clumsy Go implementation of
// private fields that's accessible by friend class
type PageInstance struct {
	Data     [common.PageSize]byte
	PageID   common.PageID
	PinCount int
	Dirty    bool
	RWLatch  common.ReaderWriterLatch
}

func NewPage() Page {
	p := &PageInstance{
		Data:     [common.PageSize]byte{},
		PageID:   common.InvalidPageID,
		PinCount: 0,
		Dirty:    false,
		RWLatch:  common.NewRWLatch(),
	}

	p.resetMemory()

	return p
}

func (p *PageInstance) resetMemory() {
	for i := range p.Data {
		p.Data[i] = offsetPageStart
	}
}

func (p *PageInstance) GetData() []byte {
	return p.Data[:]
}

func (p *PageInstance) GetPageID() common.PageID {
	return p.PageID
}

func (p *PageInstance) GetPinCount() int {
	return p.PinCount
}

func (p *PageInstance) IsDirty() bool {
	return p.Dirty
}

func (p *PageInstance) GetLSN() common.LSN {
	data := p.GetData()
	ptr := unsafe.Pointer(&data[offsetLSN])
	return *(*common.LSN)(ptr)
}

func (p *PageInstance) SetLSN(lsn common.LSN) {
	data := p.GetData()
	ptr := unsafe.Pointer(&data[offsetLSN])
	*(*common.LSN)(ptr) = lsn
}

func (p *PageInstance) WLatch() {
	p.RWLatch.WLock()
}

func (p *PageInstance) WUnlatch() {
	p.RWLatch.WUnlock()
}

func (p *PageInstance) RLatch() {
	p.RWLatch.RLock()
}

func (p *PageInstance) RUnlatch() {
	p.RWLatch.RUnlock()
}
