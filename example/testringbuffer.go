package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/wowqhb/ringbuffer/ringbuffer"

	//"strconv"
	"strconv"
)

func main() {
	rBuffer, err := ringbuffer.NewRingBuffer[[]byte](int64(32))
	fmt.Println(" ringbuffer.NewRingBuffer(int64(32)):", err)
	fmt.Println(rBuffer.GetCurrentReadIndex())
	fmt.Println(rBuffer.GetCurrentWriteIndex())
	bytes := make([]byte, 20)
	bytes[1] = byte(20)
	fmt.Println(&bytes)
	ok := rBuffer.WriteBuffer(&bytes)
	fmt.Println(ok)
	retP, ok := rBuffer.ReadBuffer()
	fmt.Println(ok)
	fmt.Println(retP)
	bytes[0] = byte(19)
	fmt.Println(bytes)
	fmt.Println(retP)
	go writeGoroutine(rBuffer)
	go readGoroutine(rBuffer)
	time.Sleep(60 * time.Second)

}

func readGoroutine(rBuffer *ringbuffer.RingBuffer[[]byte]) {
	for {
		retP, ok := rBuffer.ReadBuffer()
		if ok {
			if retP != nil {
				fmt.Println(rBuffer.GetCurrentReadIndex()-1, "::READ::", *retP, " =>> ", ok)
			}

		}
		time.Sleep(1 * time.Millisecond)
	}
}

func writeGoroutine(rBuffer *ringbuffer.RingBuffer[[]byte]) {
	for {
		time_ := strconv.FormatInt(rBuffer.GetCurrentWriteIndex()+int64(10000), 10)
		bytes := bytes.NewBufferString(time_).Bytes()
		ok := rBuffer.WriteBuffer(&bytes)
		windex := rBuffer.GetCurrentWriteIndex()
		if ok {
			fmt.Println(windex, "::WRITE::", bytes, " =>> ", ok)
		}
		time.Sleep(1 * time.Millisecond)
	}
}
