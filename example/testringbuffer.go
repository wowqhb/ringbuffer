package main

import (
	"github.com/wowqhb/ringbuffer/ringbuffer"
	"fmt"
)

func main() {
	rbuffer := ringbuffer.RingBuffer{}
	rbuffer.RingBufferInit(int64(200))
	fmt.Println(rbuffer.GetCurrentReadIndex())
	fmt.Println(rbuffer.GetCurrentWriteIndex())
}
