package ringbuffer

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

type ArrayStruct struct {
	currentTime int64
	p           []byte
	realLen     int64
	maxLen      int64
	arrayPool   *ArrayPool
}

func newArrayStruct(ap *ArrayPool) (*ArrayStruct, error) {
	return &ArrayStruct{
		currentTime: time.Now().Unix(),
		p:           make([]byte, 1024),
		realLen:     int64(0),
		maxLen:      int64(1024),
		arrayPool:   ap,
	}, nil
}

func (this *ArrayStruct) flashTime() {
	this.currentTime = time.Now().Unix()
}

func (this *ArrayStruct) ReBackToPool() {
	this.arrayPool.pool.PushFront(this)
}

type ArrayPool struct {
	max          int64
	pool         *list.List
	currentTotal int64
	lock         sync.Cond
}

func NewArrayPool(size int64) (*ArrayPool, error) {
	if int64(0) == size {
		return nil, error.Error("ERROR:NewArrayPool falure")
	}
	ap := &ArrayPool{
		pool:         list.New(),
		currentTotal: int64(0),
		max:          int64(size),
		lock:         sync.NewCond(new(sync.Mutex)),
	}

	return ap, nil
}

func (this *ArrayPool) GetMax() int64 {
	return this.max
}

func (this *ArrayPool) GetCurrentTotal() int64 {
	return atomic.LoadInt64(&this.currentTotal)
}

func (this *ArrayPool) addCurrentTotal() bool {
	if this.max <= this.GetCurrentTotal() {
		return false
	}
	atomic.AddInt64(&this.currentTotal, int64(1))
	return true
}

func (this *ArrayPool) getArrayStruct() (*ArrayStruct, error) {
	this.lock.L.Lock()
	defer this.lock.L.Unlock()
	if this.max >= this.GetCurrentTotal() {
		//存在
		if this.pool.Len() > 0 {
			f := this.pool.Front()
			this.pool.Remove(f)
			ArrayStruct{}(f.Value).flashTime()
			return f.Value, nil
		}
		//不存在，则创建
		if this.GetCurrentTotal() == 0 {
			as, err := newArrayStruct(this)
			if err == nil {
				if this.addCurrentTotal() {
					return &as, nil
				}
			}
		}
	}

	return nil, error.Error("ERROR:getArrayStruct falure")
}

func (this *ArrayPool) Cleaner() {
	if this == nil {
		return
	}
	this.lock.L.Lock()
	defer this.lock.L.Unlock()
	_tmp := list.New()
	if this.pool.Len() > 0 {
		for _, v := range this.pool {
			if v != nil {
				as := ArrayStruct{}(v)
				//时间差5分钟
				if time.Now().Unix()-as.currentTime > int64(300000) {
					_tmp.PushFront(v)
				}
			}
		}
		for i := 0; i < _tmp.Len(); i++ {
			v := _tmp.Front()
			this.pool.Remove(v)
			_tmp.Remove(v)
		}
	}
}
