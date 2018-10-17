package main

import (
	"github.com/sunguoguo/memory"
	"github.com/sunguoguo/memory/core"
	"github.com/sunguoguo/memory/example/app"
	spiders2 "github.com/sunguoguo/memory/example/app/spiders"
)

/**

爬虫启动

*/

func run() {

	Spiders := []memory.MMSpider{

		new(spiders2.ExampleSpider),
	}

	Pipelines := []memory.MMPipeline{
		//new(app.PipelineShop),
	}

	Signal := new(app.Signal)

	engine := core.MMEngineSington()
	engine.Run(Spiders, Pipelines, Signal)

}

func main() {

	run()

}
