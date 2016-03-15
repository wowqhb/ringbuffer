package main

import (
	"bytes"
	"fmt"
	"github.com/wowqhb/ringbuffer/ringbuffer"
	"time"
	//"strconv"
	"strconv"
)

func main() {
	rbuffer, err := ringbuffer.NewRingBuffer(int64(32))
	fmt.Println(" ringbuffer.NewRingBuffer(int64(32)):", err)
	fmt.Println(rbuffer.GetCurrentReadIndex())
	fmt.Println(rbuffer.GetCurrentWriteIndex())
	/*bytes := make([]byte, 20)
	bytes[1] = byte(20)
	fmt.Println(&bytes)
	ok := rbuffer.WriteBuffer(&bytes)
	fmt.Println(ok)
	retP, ok := rbuffer.ReadBuffer()
	fmt.Println(ok)
	fmt.Println(retP)
	bytes[0] = byte(19)
	fmt.Println(bytes)
	fmt.Println(retP)*/
	go writegoroutine(rbuffer)
	go readgoroutine(rbuffer)
	time.Sleep(600 * time.Second)

}

func readgoroutine(rbuffer *ringbuffer.RingBuffer) {
	for {
		retP, ok := rbuffer.ReadBuffer()
		if ok {
			if retP != nil {
				fmt.Println(rbuffer.GetCurrentReadIndex()-1, "::READ::", *retP, " =>> ", ok)
			}

		}
	}
}

func writegoroutine(rbuffer *ringbuffer.RingBuffer) {
	for {
		time_ := strconv.FormatInt(rbuffer.GetCurrentWriteIndex()+int64(10000), 10)
		bytes := bytes.NewBufferString(time_).Bytes()
		ok := rbuffer.WriteBuffer(&bytes)
		windex := rbuffer.GetCurrentWriteIndex()
		if ok {
			fmt.Println(windex, "::WRITE::", bytes, " =>> ", ok)
		}
	}
}
