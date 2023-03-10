package ringbuffer

import (
	"errors"
	"sync"
	"sync/atomic"
)

type RingBuffer[T any] struct {
	readIndex  int64      //读序号
	writeIndex int64      //写序号
	buf        []*T       //环形buffer指针数组
	bufSize    int64      //初始化环形buffer指针数组大小
	mask       int64      //初始化环形buffer指针数组大小
	pCond      *sync.Cond //生产者
	cCond      *sync.Cond //消费者
	done       int64      //is done? 1=done; 0=doing
}

func powerOfTwo64(n int64) bool {
	return n != 0 && (n&(n-1)) == 0
}

/*
*
初始化ringbuffer
参数bufferSize：初始化环形buffer指针数组大小
*/
func NewRingBuffer[T any](size int64) (*RingBuffer[T], error) {
	if !powerOfTwo64(size) {
		return nil, errors.New("Size必须为2的幂次方")
	}
	buffer := RingBuffer[T]{
		readIndex:  int64(0),
		writeIndex: int64(0),
		buf:        make([]*T, size),
		bufSize:    size,
		mask:       size - int64(1),
		pCond:      sync.NewCond(new(sync.Mutex)),
		cCond:      sync.NewCond(new(sync.Mutex)),
		done:       int64(0),
	}
	for i := int64(0); i < size; i++ {
		buffer.buf[i] = nil
	}
	return &buffer, nil
}

/*
*
获取当前读序号
*/
func (t *RingBuffer[T]) GetCurrentReadIndex() int64 {
	return atomic.LoadInt64(&t.readIndex)
}

/*
*
获取当前写序号
*/
func (t *RingBuffer[T]) GetCurrentWriteIndex() int64 {
	return atomic.LoadInt64(&t.writeIndex)
}

/*
*
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
*/
func (t *RingBuffer[T]) ReadBuffer() (p *T, ok bool) {
	t.cCond.L.Lock()
	defer func() {
		t.cCond.Signal()
		t.cCond.L.Unlock()
	}()
	ok = false
	p = nil
	readIndex := t.GetCurrentReadIndex()
	for {
		if t.isDone() {
			return nil, false
		}
		writeIndex := t.GetCurrentWriteIndex()
		if readIndex >= writeIndex {
			t.cCond.Signal()
			t.cCond.Wait()
		} else {
			break
		}
	}
	index := readIndex & t.mask //替代求模
	p = t.buf[index]
	t.buf[index] = nil
	atomic.AddInt64(&t.readIndex, int64(1))
	if p != nil {
		ok = true
	}
	return p, ok
}

/*
*
写入ringbuffer指针，以及将写序号加1
*/
func (t *RingBuffer[T]) WriteBuffer(in *T) (ok bool) {
	t.cCond.L.Lock()
	defer func() {
		t.cCond.Signal()
		t.cCond.L.Unlock()
	}()
	ok = false
	writeIndex := t.GetCurrentWriteIndex()
	for {
		if t.isDone() {
			return false
		}
		readIndex := t.GetCurrentReadIndex()
		if writeIndex >= readIndex && writeIndex-readIndex >= t.bufSize {
			t.cCond.Signal()
			t.cCond.Wait()
		} else {
			break
		}
	}
	index := writeIndex & t.mask //替代求模
	t.buf[index] = in
	atomic.AddInt64(&t.writeIndex, int64(1))
	ok = true
	return ok
}

func (t *RingBuffer[T]) Close() error {
	atomic.StoreInt64(&t.done, 1)

	t.cCond.L.Lock()
	t.cCond.Signal()
	t.cCond.L.Unlock()

	t.cCond.L.Lock()
	t.cCond.Signal()
	t.cCond.L.Unlock()

	return nil
}

func (t *RingBuffer[T]) isDone() bool {
	return atomic.LoadInt64(&t.done) == 1
}
