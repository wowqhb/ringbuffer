package ringbuffer

import (
	"sync/atomic"
)

type RingBuffer struct {
	readIndex  uint64    //读序号
	writeIndex uint64    //写序号
	ringBuffer []*[]byte //环形buffer指针数组
	bufferSize uint64    //初始化环形buffer指针数组大小
	k          uint64
}

/**
初始化ringbuffer
参数bufferSize：初始化环形buffer指针数组大小
 */
func (buffer *RingBuffer)RingBufferInit(k uint64) {
	buffer.readIndex = 0
	buffer.writeIndex = 0
	buffer.bufferSize = uint64(1) << k
	buffer.k = k
	buffer.ringBuffer = make([]*[]byte, buffer.bufferSize)
}

/**
获取当前读序号
 */
func (buffer *RingBuffer)GetCurrentReadIndex() (uint64) {
	return atomic.LoadInt64(&buffer.readIndex)
}
/**
获取当前写序号
 */
func (buffer *RingBuffer)GetCurrentWriteIndex() (uint64) {
	return atomic.LoadInt64(&buffer.writeIndex)
}

/**
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
 */
func (buffer *RingBuffer)ReadBuffer() (p *[]byte, ok bool) {
	ok = true
	p = nil

	readIndex := atomic.LoadInt64(&buffer.readIndex)
	writeIndex := atomic.LoadInt64(&buffer.writeIndex)
	switch  {
	case readIndex >= writeIndex:
		ok = false
	case writeIndex - readIndex > buffer.bufferSize:
		ok = false
	default:
		//index := buffer.readIndex % buffer.bufferSize
		index := readIndex & ((1 << buffer.k) - 1)
		p = buffer.ringBuffer[index]
		buffer.ringBuffer[index] = nil
		atomic.AddUint64(&buffer.readIndex, 1)
		if p == nil {
			ok = false
		}
	}
	return p, ok
}

/**
写入ringbuffer指针，以及将写序号加1
 */
func (buffer *RingBuffer)WriteBuffer(in *[]byte) (ok bool) {
	ok = true

	readIndex := atomic.LoadInt64(&buffer.readIndex)
	writeIndex := atomic.LoadInt64(&buffer.writeIndex)
	switch  {
	case writeIndex - readIndex < 0:
		ok = false
	default:
		//index := buffer.writeIndex % buffer.bufferSize
		index := writeIndex & ((1 << buffer.k) - 1)
		if buffer.ringBuffer[index] == nil {
			buffer.ringBuffer[index] = in
			atomic.AddUint64(&buffer.writeIndex, 1)
		}else {
			ok = false
		}
	}
	return ok
}