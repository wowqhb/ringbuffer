package ringbuffer

import (
	"sync/atomic"
)

type Ringbuffer struct {
	readIndex  uint64    //读序号
	writeIndex uint64    //写序号
	ringBuffer []*[]byte //环形buffer指针数组
	bufferSize uint64    //初始化环形buffer指针数组大小
}

/**
初始化ringbuffer
参数bufferSize：初始化环形buffer指针数组大小
 */
func (buffer *Ringbuffer)RingBufferInit(bufferSize uint64) {
	buffer.readIndex = 0
	buffer.writeIndex = 0
	buffer.bufferSize = bufferSize
	buffer.ringBuffer = make([]*[]byte, buffer.bufferSize)
}

/**
获取当前读序号
 */
func (buffer *Ringbuffer)GetCurrentReadIndex() (uint64) {
	return buffer.readIndex
}
/**
获取当前写序号
 */
func (buffer *Ringbuffer)GetCurrentWriteIndex() (uint64) {
	return buffer.writeIndex
}

/**
读取ringbuffer指定的buffer指针，返回该指针并清空ringbuffer该位置存在的指针内容，以及将读序号加1
 */
func (buffer *Ringbuffer)ReadBuffer() (p *[]byte, ok error) {
	ok = true
	p = nil
	switch  {
	case buffer.readIndex > buffer.writeIndex:
		ok = false
	case buffer.writeIndex - buffer.readIndex > buffer.bufferSize:
		ok = false
	default:
		p = &(buffer.ringBuffer[buffer.readIndex])
		buffer.ringBuffer[buffer.readIndex] = nil
		atomic.AddUint64(&buffer.readIndex, 1)
		if p == nil {
			p = false
		}
	}
	return p, ok
}

/**
写入ringbuffer指针，以及将写序号加1
 */
func (buffer *Ringbuffer)WriteBuffer(in *[]byte) (ok error) {
	ok = true
	switch  {
	case buffer.writeIndex - buffer.readIndex >= buffer.bufferSize:
		ok = false
	default:
		index := atomic.AddUint64(&buffer.writeIndex, 1)
		if buffer.ringBuffer[index] == nil {
			buffer.ringBuffer[index] = in
		}else {
			ok = false
		}
	}
	return ok
}