package ringbuffer

import (
	"fmt"
	"sync/atomic"
	"time"
)

type RingBuffer struct {
	buf  chan *ArrayStruct //环形buffer指针数组
	done int64             //is done? 1=done; 0=doing
	pool *ArrayPool
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
		buf:  make(chan *ArrayStruct, size),
		done: int64(0),
	}
	return &buffer, nil
}

/**
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
*/
func (this *RingBuffer) ReadBuffer() (*ArrayStruct, bool) {
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
	as, err := this.pool.getArrayStruct()
	if err != nil {
		return false
	}
	_len := int64(len(in))
	if _len > as.maxLen {
		as.maxLen = _len
		as.p = in
		as.realLen = _len
	} else {
		as.realLen = _len
		copy(as.p[0:], in[0:])
	}
	select {
	case this.buf <- as:
		return true
	}
	return false
}

func (this *RingBuffer) Close() {
	atomic.StoreInt64(&this.done, 1)
	close(this.buf)
}

func (this *RingBuffer) isDone() bool {
	if atomic.LoadInt64(&this.done) == 1 {
		return true
	}

	return false
}

func (this *RingBuffer) Cleaner() {
	for !this.isDone() {
		this.pool.Cleaner()
		fmt.Println("Cleaner running")
		time.Sleep(10 * time.Minute)
	}
}
