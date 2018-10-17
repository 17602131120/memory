package spiders

import "github.com/sunguoguo/memory/core"

type ExampleSpider struct{}

func (this *ExampleSpider) Pause() {

}

/**
 * 阻塞生成种子
 */

func (this *ExampleSpider) StartSeed(chum chan int, spiderName string, startSeed bool) {

	if startSeed {
		//如果队列为空，则只读取一次

	}

	//spider 执行种子结束调用的方法
	core.MMEngineSington().CloseSpider(spiderName)
	chum <- 1

}
