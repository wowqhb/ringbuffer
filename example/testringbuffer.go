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
	bytes := make([]byte, 20)
	fmt.Println(&bytes)
	ok := rbuffer.WriteBuffer(&bytes)
	fmt.Println(ok)
	retP, ok := rbuffer.ReadBuffer()
	fmt.Println(ok)
	fmt.Println(retP)
	fmt.Println(&retP)

}
