package memory

type MMSignal interface {
	ReceiptMsg(msg string)                      //接收引擎不定期信息
	ReceiptEngineLog(engineLog MemoryEngineLog) //定时接收引擎反馈的日志信息
	Open()                                      //开启引擎信号
	Close()                                     //关闭引擎信号
}
