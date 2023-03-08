package internal

import "sync"

var BufferPool = NewPool()

type Pool struct {
	p1 *sync.Pool
	p2 *sync.Pool
	p3 *sync.Pool
}

// For plugin ebpfdriver as an example, most of the data sizes are within
// 2048, drop size over 4096 as default
func NewPool() Pool {
	return Pool{
		p1: &sync.Pool{New: func() interface{} { return make([]byte, 2<<10) }},
		p2: &sync.Pool{New: func() interface{} { return make([]byte, 4<<10) }},
		p3: &sync.Pool{New: func() interface{} { return make([]byte, 8<<10) }},
	}
}

func (p *Pool) Get(size int64) []byte {
	if size <= 2<<10 {
		return p.p1.Get().([]byte)[:size-1]
	} else if size <= 4<<10 {
		return p.p2.Get().([]byte)[:size-1]
	} else if size <= 8<<10 {
		return p.p3.Get().([]byte)[:size-1]
	}
	return make([]byte, size)
}

func (p *Pool) Put(b []byte) {
	if len(b) <= 2<<10 {
		p.p1.Put(b)
		return
	}
	if len(b) <= 4<<10 {
		p.p2.Put(b)
		return
	}
	if len(b) <= 8<<10 {
		p.p3.Put(b)
		return
	}
	// deprecate []byte over 4096, let GC collects them
}
