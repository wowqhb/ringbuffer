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
				_i, err := strconv.Atoi(bytes.NewBuffer(retP.GetBytes()).String())
				if err != nil {
					ok = false
				} else {
					ok = i == _i
				}
				fmt.Println(i, "::READ::", retP.GetBytes(), " =>> ", ok)
				retP.Destroy()
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
		bs := rbuffer.CreateBufferStruct()
		bs.Check(len(_bytes))
		copy(bs.GetBytes()[0:], _bytes[0:])
		ok := rbuffer.WriteBuffer(bs)
		if ok {
			fmt.Println(i, "::WRITE::", bytes.NewBuffer(_bytes).String(), " =>> ", ok)
			i++
		}
	}
}
