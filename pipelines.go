package memory

type MMPipeline interface {
	ProcessItem(item MMItem) //持久化item回调函数
	Open(spiderName string)  //开启pipeline，配置多少持久化线程，就调用多少次函数
	Close(spiderName string) //关闭pipeline，配置多少持久化线程，就调用多少次函数
}
