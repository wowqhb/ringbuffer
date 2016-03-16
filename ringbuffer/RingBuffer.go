package ringbuffer

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type RingBuffer struct {
	readIndex  int64      //读序号
	writeIndex int64      //写序号
	buf        []*[]byte  //环形buffer指针数组
	bufSize    int64      //初始化环形buffer指针数组大小
	mask       int64      //初始化环形buffer指针数组大小
	pcond      *sync.Cond //生产者
	ccond      *sync.Cond //消费者
	done       int64      //is done? 1=done; 0=doing
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
		readIndex:  int64(0),
		writeIndex: int64(0),
		buf:        make([]*[]byte, size),
		bufSize:    size,
		mask:       size - int64(1),
		pcond:      sync.NewCond(new(sync.Mutex)),
		ccond:      sync.NewCond(new(sync.Mutex)),
		done:       int64(0),
	}
	for i := int64(0); i < size; i++ {
		buffer.buf[i] = nil
	}
	return &buffer, nil
}

/**
获取当前读序号
*/
func (this *RingBuffer) GetCurrentReadIndex() int64 {
	return atomic.LoadInt64(&this.readIndex)
}

/**
获取当前写序号
*/
func (this *RingBuffer) GetCurrentWriteIndex() int64 {
	return atomic.LoadInt64(&this.writeIndex)
}

/**
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
*/
func (this *RingBuffer) ReadBuffer() (p *[]byte, ok bool) {
	this.ccond.L.Lock()
	defer func() {
		this.pcond.Broadcast()
		this.ccond.L.Unlock()
	}()
	ok = false
	p = nil
	readIndex := this.GetCurrentReadIndex()
	writeIndex := this.GetCurrentWriteIndex()
	for {
		if this.isDone() {
			return nil, false
		}
		writeIndex = this.GetCurrentWriteIndex()
		if readIndex >= writeIndex {
			//fmt.Println("read wait")
			this.pcond.Broadcast()
			this.ccond.Wait()
		} else {
			break
		}
		//time.Sleep(1 * time.Millisecond)
	}
	index := readIndex & this.mask //替代求模
	p = this.buf[index]
	this.buf[index] = nil
	atomic.AddInt64(&this.readIndex, int64(1))
	if p != nil {
		ok = true
	}
	return p, ok
}

/**
写入ringbuffer指针，以及将写序号加1
*/
func (this *RingBuffer) WriteBuffer(in *[]byte) (ok bool) {
	this.pcond.L.Lock()
	defer func() {
		this.ccond.Broadcast()
		this.pcond.L.Unlock()
	}()
	ok = false
	readIndex := this.GetCurrentReadIndex()
	writeIndex := this.GetCurrentWriteIndex()
	for {
		if this.isDone() {
			return false
		}
		readIndex = this.GetCurrentReadIndex()
		if writeIndex >= readIndex && writeIndex-readIndex >= this.bufSize {
			//fmt.Println("write wait")
			this.ccond.Broadcast()
			this.pcond.Wait()
			//time.Sleep(1 * time.Millisecond)
		} else {
			break
		}
		//time.Sleep(1 * time.Millisecond)

	}
	index := writeIndex & this.mask //替代求模
	this.buf[index] = in
	atomic.AddInt64(&this.writeIndex, int64(1))
	ok = true
	return ok
}

func (this *RingBuffer) Close() error {
	atomic.StoreInt64(&this.done, 1)

	this.pcond.L.Lock()
	this.ccond.Broadcast()
	this.pcond.L.Unlock()

	this.ccond.L.Lock()
	this.pcond.Broadcast()
	this.ccond.L.Unlock()

	return nil
}

func (this *RingBuffer) isDone() bool {
	if atomic.LoadInt64(&this.done) == 1 {
		return true
	}

	return false
}
