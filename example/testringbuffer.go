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

	time.Sleep(60 * time.Second)

}

func readgoroutine(rbuffer *ringbuffer.RingBuffer) {
	i := 0
	for {
		retP, ok := rbuffer.ReadBuffer()
		if ok {
			if retP != nil {
				i++
				_i, err := strconv.Atoi(bytes.NewBuffer(retP).String())
				if err != nil {
					ok = false
				} else {
					ok = i == _i
				}
				fmt.Println(i, "::READ::", bytes.NewBuffer(retP.GetBytes()).String(), " =>> ", ok)
				retP.ReBackToPool()
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
		ok := rbuffer.WriteBuffer(_bytes)
		if ok {
			fmt.Println(i, "::WRITE::", bytes.NewBuffer(_bytes).String(), " =>> ", ok)
			i++
		}
	}
}
