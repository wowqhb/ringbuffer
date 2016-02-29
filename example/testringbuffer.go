package main

import (
	"github.com/wowqhb/ringbuffer/ringbuffer"
	"fmt"
	"time"
	"bytes"
	"strconv"
)

func main() {
	rbuffer := ringbuffer.RingBuffer{}
	rbuffer.RingBufferInit(3)
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
	go writegoroutine(&rbuffer)
	go readgoroutine(&rbuffer)

	time.Sleep(60 * time.Second)

}

func readgoroutine(rbuffer  *ringbuffer.RingBuffer) {
	for {
		fmt.Println(strconv.FormatInt(rbuffer.GetCurrentWriteIndex(),10) + "::::" +strconv.FormatInt(rbuffer.GetCurrentReadIndex(),10))
		retP, ok := rbuffer.ReadBuffer()
		if ok {
			fmt.Print("read::")
			fmt.Print(retP)
			fmt.Print(" =>> ")
			fmt.Println(ok)
		}else {
			fmt.Print("read::nil =>> ")
			fmt.Println(ok)
		}
		//time.Sleep(10 * time.Millisecond)
	}
}

func writegoroutine(rbuffer *ringbuffer.RingBuffer) {
	for {
		fmt.Println(strconv.FormatInt(rbuffer.GetCurrentWriteIndex(),10) + "::::" +strconv.FormatInt(rbuffer.GetCurrentReadIndex(),10))
		time_ := time.Now().String()
		bytes := bytes.NewBufferString(time_).Bytes()
		ok := rbuffer.WriteBuffer(&bytes)
		fmt.Print("write::" + time_ + " =>> ")
		fmt.Println(ok)
		//time.Sleep(10 * time.Millisecond)
	}
}