(1) time.Sleep(1*time.Microsecond)//读取或写入失败后休眠1微秒
(2) time.Sleep(1 * time.Millisecond) //读取或写入失败后休眠1毫秒
(3) 无休眠
(4) runtime.Gosched()//读取或写入失败后让出时间片

性能：(4) > (2) > (1) > (3)