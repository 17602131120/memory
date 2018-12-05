package core

import (
	"fmt"
	"github.com/sunguoguo/memory"
	"github.com/sunguoguo/memory/utils"
	"sync"
)

type mmScheduler struct {
	mmRedis *utils.MMRedis
	mmMongo *utils.MMMongo
}

var scheduler *mmScheduler
var schedulerOnce sync.Once

func MMSchedulerSington() *mmScheduler {

	schedulerOnce.Do(func() {

		scheduler = new(mmScheduler)
		scheduler.mmRedis = new(utils.MMRedis)
		scheduler.mmMongo = new(utils.MMMongo)

	})

	return scheduler

}

func (this *mmScheduler) RequestPull(key string) (*memory.MemoryRequest, error) {

	return this.mmRedis.RequestPull(key)

}
func (this *mmScheduler) RequestPush(leftPush bool, retry bool, key string, request memory.MemoryRequest) {

	if retry { //重试的任务忽略去重

		//重试任务超过最高次数，直接舍弃任务
		if request.RequestNum > 1000 {
			return

		}

		//重试任务超过指定次数，则改变入队方向
		if request.RequestNum > 5 {

			leftPush = false
		}

		request.RequestNum += 1

	} else { //去重处理

		if request.DontFilter == false {
			if this.mmRedis.FilterRepeat(fmt.Sprintf("%sfilter", key), request.Url) {
				//已经存在，该任务已经重复
				return
			}
		}

	}

	this.mmRedis.RequestPush(leftPush, key, request)

}
func (this *mmScheduler) LLen(key string) int {

	return this.mmRedis.LLen(key)

}
func (this *mmScheduler) ItemPush(key string, item memory.MMItem) {

	this.mmRedis.ItemPush(key, item)

}
func (this *mmScheduler) ItemPull(key string) memory.MMItem {

	return this.mmRedis.ItemPull(key)

}
func (this *mmScheduler) SeedPull(key string) (bool, string) {

	return this.mmRedis.SeedPull(key)
}
func (this *mmScheduler) SeedPush(key, seed string) {

	this.mmRedis.SeedPush(key, seed)
}
