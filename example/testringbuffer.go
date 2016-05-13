package main

import (
	"bytes"
	"fmt"
	"github.com/wowqhb/ringbuffer/ringbuffer"
	"time"
)

func main() {
	rbuffer, err := ringbuffer.NewRingBuffer(int64(32))
	fmt.Println(" ringbuffer.NewRingBuffer(int64(32)):", err)
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

	go rbuffer.clea
	time.Sleep(60 * time.Second)

}

func readgoroutine(rbuffer *ringbuffer.RingBuffer) {
	i := int64(0)
	for {
		retP, ok := rbuffer.ReadBuffer()
		if ok {
			if retP != nil {
				i++
				fmt.Println(i, "::READ::", retP, " =>> ", ok)
			}

		}
	}
}

func writegoroutine(rbuffer *ringbuffer.RingBuffer) {
	i := int64(0)
	for {
		time_ := time.Now().String()
		bytes := bytes.NewBufferString(time_).Bytes()
		ok := rbuffer.WriteBuffer(bytes)
		if ok {
			i++
			fmt.Println(i, "::WRITE::", bytes, " =>> ", ok)
		}
	}
}
