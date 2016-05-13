package ringbuffer

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type RingBuffer struct {
	buf  chan []byte //环形buffer指针数组
	done int64       //is done? 1=done; 0=doing
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
		buf:  make(chan []byte, size),
		done: int64(0),
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 1)
			},
		},
	}
	return &buffer, nil
}

/**
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
*/
func (this *RingBuffer) ReadBuffer() ([]byte, bool) {
	select {
	case p, ok := <-this.buf:
		return p, ok
	}
	return nil, false
}

/**
写入ringbuffer指针，以及将写序号加1
*/
func (this *RingBuffer) WriteBuffer(in []byte) bool {
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

func (this *RingBuffer) GetBytes(_len int) []byte {
	bs := this.pool.Get().([]byte)
	if _len > len(bs) {
		bs = make([]byte, _len)
	}
	return bs
}

func (this *RingBuffer) DestoryBytes(in []byte) {
	this.pool.Put(in)
}
