package memory

type MMSpider interface {
	StartSeed(chum chan int, spiderName string, startSeed bool)

	Pause()
}
