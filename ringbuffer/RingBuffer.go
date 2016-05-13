package ringbuffer

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type BufferStruct struct {
	realLen int
	p       []byte
	pool    sync.Pool
}

func (this *BufferStruct) GetBytes() []byte {
	return this.p[0:this.realLen]
}

func (this *BufferStruct) Check(len int) {
	if len > len(this.p) {
		this.p = make([]byte, len)
	}
}
func (this *BufferStruct) Destroy() {
	this.pool.Put(this)
}

type RingBuffer struct {
	buf  chan *BufferStruct //环形buffer指针数组
	done int64              //is done? 1=done; 0=doing
	pool *sync.Pool
}

func powerOfTwo64(n int64) bool {
	return n != 0 && (n&(n-1)) == 0
}

/**
初始化ringbuffer
参数bufferSize：初始化环形buffer指针数组大小
*/
func NewRingBuffer(size int64) (*RingBuffer, error) {
	if !powerOfTwo64(size) {
		return nil, fmt.Errorf("This size is not able to used")
	}
	buffer := RingBuffer{
		buf:  make(chan *BufferStruct, size),
		done: int64(0),
		pool: &sync.Pool{
			New: func() interface{} {
				return &BufferStruct{
					realLen: 0,
					p:       make([]byte, 8192),
					pool:    nil,
				}
			},
		},
	}
	return &buffer, nil
}

/**
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
*/
func (this *RingBuffer) ReadBuffer() (*BufferStruct, bool) {
	select {
	case p, ok := <-this.buf:
		return p, ok
	}
	return nil, false
}

/**
写入ringbuffer指针，以及将写序号加1
*/
func (this *RingBuffer) WriteBuffer(in *BufferStruct) bool {
	select {
	case this.buf <- in:
		return true
	}
	return false
}

func (this *RingBuffer) Close() {
	atomic.StoreInt64(&this.done, 1)
	this.pool = nil
	close(this.buf)
}

func (this *RingBuffer) isDone() bool {
	if atomic.LoadInt64(&this.done) == 1 {
		return true
	}

	return false
}

func (this *RingBuffer) CreateBufferStruct() *BufferStruct {
	bs := this.pool.Get().(*BufferStruct)
	bs.pool = this.pool
	return bs
}
