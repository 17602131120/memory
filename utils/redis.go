package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/memory"
	"github.com/memory/setting"
	"sync"
	"time"
)

var redisClient *redis.Client
var redisOnce sync.Once

type MMRedis struct {
}

func (this *MMRedis) getClient() *redis.Client {

	redisOnce.Do(func() {

		settings := setting.MMSettingsSington()

	CONNECT:
		{
			redisClient = redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", settings.Redis.Host, settings.Redis.Port),
				Password: settings.Redis.Password, // no password set
				DB:       settings.Redis.Db,       // use default DB
			})
			_, err := redisClient.Ping().Result()

			if err != nil {
				//连接失败
				logger := setting.MMSettingsSington().MMLogger
				logger.Printf("redis 连接无法建立，休眠5s重试...")
				time.Sleep(time.Second * 5)

				goto CONNECT
			}
		}

	})

	return redisClient
}
func (this *MMRedis) LLen(key string) int {
	client := this.getClient()
	res := client.LLen(key)
	return int(res.Val())

}

//从指定队列获取请求结构体
func (this *MMRedis) RequestPull(key string) (memory.MemoryRequest, error) {

	client := this.getClient()
	var request memory.MemoryRequest

	err := json.Unmarshal([]byte(client.LPop(key).Val()), &request)

	if err != nil {
		//未获取到数据，失败

		setting.MMSettingsSington().MMLogger.Printf("RequestPull failed: %s （key=%s）", err, key)

		return request, err
	} else {
		//成功获取数据,结束

		return request, nil

	}

}

func (this *MMRedis) ItemPull(key string) memory.MMItem {
	client := this.getClient()

	var item memory.MMItem

	result := client.LPop(key).Val()

	if "" != result {
		err := json.Unmarshal([]byte(result), &item)

		if err != nil {
			//未获取到数据，失败
			logger := setting.MMSettingsSington().MMLogger
			logger.Printf("ItemPull failed: %s （key=%s）（item=%s）", err, key, item)

			return nil
		} else {
			//成功获取数据,结束
			return item
		}

	} else {
		return nil
	}

}
func (this *MMRedis) ItemPush(key string, item memory.MMItem) {
	client := this.getClient()

	b, err := json.Marshal(item)
	if err != nil {
		//失败
		logger := setting.MMSettingsSington().MMLogger
		logger.Printf("ItemPush failed: %s （key=%s）（item=%s）", err, key, item)
	} else {
		client.RPush(key, string(b))

	}

}

//push请求结构体到指定队列，可以自定义入队方向
func (this *MMRedis) RequestPush(leftPush bool, key string, request memory.MemoryRequest) {
	client := this.getClient()

	b, err := json.Marshal(request)
	if err != nil {

		//失败
		logger := setting.MMSettingsSington().MMLogger
		logger.Printf("RequestPush failed: %s （key=%s）（request=%s）", err, key, request)

	} else {

		if leftPush {
			client.LPush(key, string(b))
		} else {
			client.RPush(key, string(b))
		}
	}
}

//指定key添加userAgent到set
func (this *MMRedis) UserAgentSadd(key string, mmUserAgent memory.MemoryUserAgent) {

	client := this.getClient()

	b, err := json.Marshal(mmUserAgent)
	if err != nil {
		//失败
		logger := setting.MMSettingsSington().MMLogger
		logger.Printf("UserAgentSadd json.Marshal failed:%s", err)
	} else {
		if mmUserAgent.Mobile {
			key = fmt.Sprintf("%sMobile", key)
		}
		client.SAdd(key, string(b))
	}

}

func (this *MMRedis) FilterRepeat(key string, val string) bool {
	client := this.getClient()

	var exist bool
	if client.SIsMember(key, val).Val() {
		//已经存在，
		exist = true

	} else {
		//不存在，插入

		client.SAdd(key, val)
		exist = false
	}

	return exist

}
func (this *MMRedis) UserAgentSrandmember(key string, mobile bool) memory.MemoryUserAgent {

	client := this.getClient()

	keyChange := key
	if mobile {
		keyChange = fmt.Sprintf("%sMobile", keyChange)
	}
	value := client.SRandMember(keyChange)

	var mmUserAgent memory.MemoryUserAgent
	err := json.Unmarshal([]byte(value.Val()), &mmUserAgent)

	if err != nil {
		logger := setting.MMSettingsSington().MMLogger
		logger.Printf("UserAgentSrandmember failed （key=%s）（err=%s）", key, err)
		time.Sleep(60 * time.Second)
		//重新获取
		return this.UserAgentSrandmember(key, mobile)
	} else {
		return mmUserAgent
	}

	return mmUserAgent

}

//指定key添加cookie到set
func (this *MMRedis) CookieSadd(key string, mmCookie memory.MemoryCookie) {

	b, err := json.Marshal(mmCookie)
	if err != nil {
		//失败
		logger := setting.MMSettingsSington().MMLogger
		logger.Printf("CookieSadd json.Marshal failed:%s", err)

	} else {
		client := this.getClient()
		client.SAdd(key, string(b))
	}

}

//指定key获取cookie
func (this *MMRedis) CookieSrandmember(key string) *memory.MemoryCookie {
	client := this.getClient()
	value := client.SRandMember(key)
	var mmCookie *memory.MemoryCookie

	err := json.Unmarshal([]byte(value.Val()), &mmCookie)

	if err != nil {
		logger := setting.MMSettingsSington().MMLogger
		logger.Printf("CookieSrandmember failed （key=%s）（err=%s）", key, err)
		time.Sleep(60 * time.Second)
		//重新获取
		return this.CookieSrandmember(key)
	} else {
		return mmCookie
	}
}

func (this *MMRedis) ProxySadd(key string, dataAll []string) {
	client := this.getClient()

	//清空原有数据
	client.Del(key)
	//添加数据
	client.SAdd(key, dataAll)

}
func (this *MMRedis) ProxySrandmember(key string) string {

	client := this.getClient()
	var proxy string

	for {
		proxy = client.SRandMember(key).Val()
		if proxy == "" {
			//获取失败,延迟一段时间，继续获取
			logger := setting.MMSettingsSington().MMLogger
			logger.Printf("ProxySrandmember failed （key=%s）")
			time.Sleep(time.Second * 60)
			continue
		} else {
			//获取成功
			break
		}
	}

	return proxy

}

//push种子到指定队列
func (this *MMRedis) SeedPush(key, seed string) {

	client := this.getClient()
	//defer client.Close()

	client.RPush(key, seed)

}

//从指定队列获取请求种子
func (this *MMRedis) SeedPull(key string) (bool, string) {
	client := this.getClient()

	//defer client.Close()

	value := client.LPop(key)
	seed := value.Val()

	var state bool = true
	if "" == seed {
		logger := setting.MMSettingsSington().MMLogger
		logger.Printf("SeedPull.未获取到种子(key=%s)", key)
		state = false
	}
	return state, seed

}
