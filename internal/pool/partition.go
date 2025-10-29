package pool

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

// Partition defines the contract implemented by pooled partitions.
type Partition interface {
	io.ReadWriter
	io.WriterAt
	Len() int64
	Bytes() []byte
	Reset()
}

// PartitionPool manages reusable partitions backed by bytes.Buffer instances.
type PartitionPool struct {
	size int
	pool sync.Pool
}

// NewPartitionPool returns a pool that provides partitions with an initial capacity.
func NewPartitionPool(size int) *PartitionPool {
	if size < 0 {
		size = 0
	}
	p := &PartitionPool{size: size}
	p.pool.New = func() any {
		part := &partition{pool: p}
		if size > 0 {
			part.buf.Grow(size)
		}
		return part
	}
	return p
}

// NewPartition acquires a partition from the provided pool.
func NewPartition(p *PartitionPool) Partition {
	if p == nil {
		return nil
	}
	return p.get()
}

func (p *PartitionPool) get() Partition {
	part := p.pool.Get().(*partition)
	part.pool = p
	part.released = false
	part.buf.Reset()
	if p.size > 0 {
		currentCap := cap(part.buf.Bytes())
		if currentCap < p.size {
			part.buf.Grow(p.size - currentCap)
		}
	}
	return part
}

type partition struct {
	pool     *PartitionPool
	buf      bytes.Buffer
	released bool
}

func (p *partition) ensureActive() error {
	if p == nil {
		return errors.New("pool: nil partition")
	}
	if p.released {
		return errors.New("pool: use after reset")
	}
	if p.pool == nil {
		return errors.New("pool: no backing pool")
	}
	return nil
}

func (p *partition) Read(b []byte) (int, error) {
	if err := p.ensureActive(); err != nil {
		return 0, err
	}
	return p.buf.Read(b)
}

func (p *partition) Write(b []byte) (int, error) {
	if err := p.ensureActive(); err != nil {
		return 0, err
	}
	return p.buf.Write(b)
}

func (p *partition) WriteAt(b []byte, off int64) (int, error) {
	if err := p.ensureActive(); err != nil {
		return 0, err
	}
	if off < 0 {
		return 0, errors.New("pool: negative offset")
	}
	if off > int64(p.buf.Len()) {
		return 0, io.ErrShortWrite
	}
	data := p.buf.Bytes()
	idx := int(off)
	if idx+len(b) > len(data) {
		return 0, io.ErrShortWrite
	}
	copy(data[idx:], b)
	return len(b), nil
}

func (p *partition) Len() int64 {
	if p == nil || p.released {
		return 0
	}
	return int64(p.buf.Len())
}

func (p *partition) Bytes() []byte {
	if p == nil || p.released {
		return nil
	}
	return p.buf.Bytes()
}

func (p *partition) Reset() {
	if p == nil || p.pool == nil || p.released {
		return
	}
	p.buf.Reset()
	p.released = true
	p.pool.pool.Put(p)
}
