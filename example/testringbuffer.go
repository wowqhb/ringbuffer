package main

import (
	"github.com/wowqhb/ringbuffer/ringbuffer"
	"fmt"
	"time"
	"bytes"
//"strconv"
	"strconv"
)

func main() {
	rbuffer := ringbuffer.RingBuffer{}
	rbuffer.RingBufferInit(uint64(8))
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
	go writegoroutine(&rbuffer)
	go readgoroutine(&rbuffer)
	time.Sleep(60 * time.Second)

}

func readgoroutine(rbuffer  *ringbuffer.RingBuffer) {
	for {
		retP, ok := rbuffer.ReadBuffer()
		if ok {
			if retP == nil {
				//fmt.Println(strconv.FormatUint(rbuffer.GetCurrentReadIndex() - 1, 10) + "::READ::nil =>> " + strconv.FormatBool(ok))
			}else {
				fmt.Println(strconv.FormatUint(rbuffer.GetCurrentReadIndex() - 1, 10) + "::READ::" + string(*retP) + " =>> " + strconv.FormatBool(ok))
			}

		}else {
			//fmt.Println(strconv.FormatUint(rbuffer.GetCurrentReadIndex(), 10) + "::READ::nil =>> " + strconv.FormatBool(ok))
		}
	}
}

func writegoroutine(rbuffer *ringbuffer.RingBuffer) {
	for {
		time_ := strconv.FormatUint(rbuffer.GetCurrentWriteIndex(), 10);
		bytes := bytes.NewBufferString(time_).Bytes()
		ok := rbuffer.WriteBuffer(&bytes)
		windex := rbuffer.GetCurrentWriteIndex()
		if ok {
			windex = rbuffer.GetCurrentWriteIndex() - 1
			fmt.Println(strconv.FormatUint(windex, 10) + "::WRITE::" + time_ + " =>> " + strconv.FormatBool(ok))
		}
		//fmt.Println(strconv.FormatUint(windex, 10) + "::WRITE::" + time_ + " =>> " + strconv.FormatBool(ok))
	}
}