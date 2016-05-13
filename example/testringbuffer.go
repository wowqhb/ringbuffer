package main

import (
	"bytes"
	"fmt"
	"github.com/wowqhb/ringbuffer/ringbuffer"
	"strconv"
	"time"
)

func main() {
	rbuffer, err := ringbuffer.NewRingBuffer(int64(32))
	fmt.Println(" ringbuffer.NewRingBuffer(int64(32)):", err)
	go writegoroutine(rbuffer)
	go readgoroutine(rbuffer)

	time.Sleep(60 * time.Second)

}

func readgoroutine(rbuffer *ringbuffer.RingBuffer) {
	i := 0
	for {
		retP, ok := rbuffer.ReadBuffer()
		if ok {
			if retP != nil {
				i++
				fmt.Println(i, "::READ::", retP, " =>> ", ok)
				rbuffer.DestoryBytes(retP)
			}

		}
	}
}

func writegoroutine(rbuffer *ringbuffer.RingBuffer) {
	i := int(1)
	for {
		time_ := strconv.Itoa(i)
		buffer := bytes.NewBufferString(time_)
		_bytes := buffer.Bytes()
		bs := rbuffer.GetBytes(len(_bytes))
		ok := rbuffer.WriteBuffer(bs)
		if ok {
			fmt.Println(i, "::WRITE::", bytes.NewBuffer(_bytes).String(), " =>> ", ok)
			i++
		}
	}
}
