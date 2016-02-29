package ringbuffer

import (
	"sync/atomic"
)

type RingBuffer struct {
	readIndex  int64     //读序号
	writeIndex int64     //写序号
	ringBuffer []*[]byte //环形buffer指针数组
	bufferSize int64     //初始化环形buffer指针数组大小
}

/**
初始化ringbuffer
参数bufferSize：初始化环形buffer指针数组大小
 */
func (buffer *RingBuffer)RingBufferInit(bufferSize int64) {
	buffer.readIndex = -1
	buffer.writeIndex = -1
	buffer.bufferSize = bufferSize
	buffer.ringBuffer = make([]*[]byte, buffer.bufferSize)
}

/**
获取当前读序号
 */
func (buffer *RingBuffer)GetCurrentReadIndex() (int64) {
	return buffer.readIndex
}
/**
获取当前写序号
 */
func (buffer *RingBuffer)GetCurrentWriteIndex() (int64) {
	return buffer.writeIndex
}

/**
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
 */
func (buffer *RingBuffer)ReadBuffer() (p *[]byte, ok bool) {
	ok = true
	p = nil
	switch  {
	case buffer.readIndex > buffer.writeIndex:
		ok = false
	case buffer.writeIndex - buffer.readIndex > buffer.bufferSize:
		ok = false
	default:
		index := buffer.readIndex % buffer.bufferSize
		p = buffer.ringBuffer[index]
		buffer.ringBuffer[index] = nil
		atomic.AddInt64(&buffer.readIndex, 1)
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
	switch  {
	case buffer.writeIndex - buffer.readIndex < 0:
		ok = false
	default:
		index := atomic.AddInt64(&buffer.writeIndex, 1) % buffer.bufferSize
		if buffer.ringBuffer[index] == nil {
			buffer.ringBuffer[index] = in
		}else {
			ok = false
		}
	}
	return ok
}