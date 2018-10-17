package app

import "github.com/sunguoguo/memory"

type Signal struct {
}

func (this *Signal) ReceiptMsg(msg string) {
	//由引擎主动调用的方法

}
func (this *Signal) ReceiptEngineLog(engineLog memory.MemoryEngineLog) {
	//定时被引擎调用的方法

}
func (this *Signal) Open() {
	//引擎开启时调用的方法

}
func (this *Signal) Close() {
	//引擎关闭时调用的方法

}
