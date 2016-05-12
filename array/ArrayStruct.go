package array

type ArrayStruct struct {
	p       []byte
	realLen int64
	maxLen  int64
}

func (this *ArrayStruct) NewArrayStruct() *ArrayStruct {
	return ArrayStruct{
		p:       make([]byte, 1024),
		realLen: int64(0),
		maxLen:  int64(1024),
	}
}
