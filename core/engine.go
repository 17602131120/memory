package core

import (
	"fmt"
	"github.com/memory"
	"github.com/memory/http"
	"github.com/memory/setting"
	"github.com/memory/utils"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

/**
爬虫引擎



*/
type mmEngine struct {
	engineLog   memory.MemoryEngineLog
	Spiders     []memory.MMSpider
	spidersName []string //和spider相同顺序的SpiderName
	Pipelines   []memory.MMPipeline
	Signal      memory.MMSignal

	closeEngine bool //默认为false 不关闭引擎 true 关闭引擎
	mmRequest   *http.HttpRequest
	mmSettings  *setting.MMSettings
	mmUtils     *utils.MMUtil
}

var engineConcurrentManage = struct { //引擎并发管理
	sync.RWMutex
	spidersState map[string]bool
}{
	spidersState: make(map[string]bool),
}

var engine *mmEngine
var engineOnce sync.Once

func MMEngineSington() *mmEngine {

	engineOnce.Do(func() {
		engine = &mmEngine{}
	})
	return engine

}

/**
关闭spider，适用多spider运行场景
*/
func (this *mmEngine) CloseSpider(SpiderName string) {

	engineConcurrentManage.Lock()
	engineConcurrentManage.spidersState[SpiderName] = false
	engineConcurrentManage.Unlock()

}

/**
检测是否关闭引擎
1：检测引擎中所有的spider是否全部关闭，如果全部关闭，允许进入第二步检测
2：检测队列情况，如果全部消费完毕，则允许关闭引擎

*/

func (this *mmEngine) _checkEngine(chum chan int) {

	scheduler = MMSchedulerSington()
	closeMaxCount := 4
	currentArray := []bool{}

	//判断是否满足关闭的条件
	judgeClose := func(array []bool) bool {

		for i := 0; i < len(array); i++ {
			if array[i] == false {
				return false
			}

		}

		return true
	}
	for {

		for _, spiderName := range this.spidersName {

			if engineConcurrentManage.spidersState[spiderName] == true {
				//蜘蛛还在运行，无法关闭引擎
				currentArray = append(currentArray, false)

			} else {
				//检测种子队列 and 检查持久化存储队列
				if scheduler.LLen(fmt.Sprintf("queue%s", spiderName)) == 0 && scheduler.LLen(fmt.Sprintf("item%s", spiderName)) == 0 {

					currentArray = append(currentArray, true)

				} else {
					currentArray = append(currentArray, false)
				}

			}

		}

		currentArrayLen := len(currentArray)

		if currentArrayLen >= closeMaxCount {
			if judgeClose(currentArray[currentArrayLen-closeMaxCount : currentArrayLen]) {
				//最近指定次数均满足关闭条件

				this.closeEngine = true
				chum <- 1

			} else {
				//不需关闭
				currentArray = currentArray[1:currentArrayLen]

			}
		}
		time.Sleep(time.Second * 5)

	}

}

func (this *mmEngine) Run(Spiders []memory.MMSpider, Pipelines []memory.MMPipeline, Signal memory.MMSignal) {

	// step1 参数初始化
	this.Signal = Signal
	this.closeEngine = false
	this.Spiders = Spiders
	this.Pipelines = Pipelines
	this.mmRequest = new(http.HttpRequest)
	this.mmSettings = setting.MMSettingsSington()
	this.mmUtils = &utils.MMUtil{}

	ConcurrentRequest := this.mmSettings.Config.ConcurrentRequest
	ConcurrentPipeline := this.mmSettings.Config.ConcurrentPipeline
	num := 1 + 1 + len(this.Spiders) + ConcurrentRequest + ConcurrentPipeline //引擎控制并发数据=蜘蛛数量+下载器+管道
	chum := make(chan int, num)

	this.engineLog = memory.MemoryEngineLog{
		MainID:    this.mmUtils.GetgoID(),
		StartTime: time.Now(),
	}
	this.mmSettings.MMLogger.Printf("Memory %s started (bot: %s) \n", memory.Version(), this.mmSettings.Config.Botname)

	for _, Spider := range this.Spiders {
		//动态为spider设置name
		t := reflect.ValueOf(Spider).Type() // *spiders.name
		nameArray := strings.Split(t.String(), ".")
		spiderName := nameArray[1]

		this.spidersName = append(this.spidersName, spiderName)
		//将spider全部放入状态map
		engineConcurrentManage.Lock()
		engineConcurrentManage.spidersState[spiderName] = true
		engineConcurrentManage.Unlock()

	}
	//启动引擎信号数量
	go this._signal(chum)
	go this._checkEngine(chum)

	// 启动
	for i := 0; i < ConcurrentRequest; i++ {

		go this._downloader(chum)
	}
	for i := 0; i < ConcurrentPipeline; i++ {
		go this._pipelineer(chum)
	}

	//启动
	scheduler = MMSchedulerSington()
	for index, Spider := range this.Spiders {
		spiderName := this.spidersName[index]
		startSeed := true
		if scheduler.LLen(fmt.Sprintf("queue%s", spiderName)) > 0 {
			//队列不为空，存在未消费完全的种子
			startSeed = false
		}

		go Spider.StartSeed(chum, spiderName, startSeed)
	}
	//统一关闭子线程
	for i := 0; i < num; i++ {
		<-chum
	}

	//引擎关闭，输出运行日志，然后退出
	this.engineLog.FinishTime = time.Now()
	this.engineLog.FinishReason = "finished"

	//this.mmSettings.MMLogger.Printf("%+v",this.engineLog)
	t := reflect.TypeOf(this.engineLog)
	v := reflect.ValueOf(this.engineLog)
	for i := 0; i < t.NumField(); i++ {
		this.mmSettings.MMLogger.Printf("%s:%v", t.Field(i).Name, v.Field(i).Interface())

	}
	//this.mmSettings.MMLogger.Printf("%+v",this.engineLog)
	this.mmSettings.MMLogger.Printf("Memory Spider closed (finished)")

}

func (this *mmEngine) _signal(chum chan int) {

	go this.Signal.Open()

	for {
		time.Sleep(time.Second * 1)
		this.Signal.ReceiptEngineLog(this.engineLog)
		time.Sleep(time.Second * 30)

		if this.closeEngine {
			this.Signal.Close()
			//关闭引擎
			chum <- 1
			return
		}

	}

}
func (this *mmEngine) _downloader(chum chan int) {
	pc, _, _, _ := runtime.Caller(0)
	scheduler = MMSchedulerSington()

	for {

		//获取所有需要被下载器消费的队列，queue+爬虫Name就是队列名
		for index, Spider := range this.Spiders {
			//spider名称
			spiderName := this.spidersName[index]
			//拼接queueKey
			queueKey := fmt.Sprintf("queue%s", spiderName)

			request, err := scheduler.RequestPull(queueKey)
			if err == nil {
				//获取到请求任务
				if request.CallbackStuct == spiderName {
					//为当前任务，找到所属蜘蛛，映射调用
					t := reflect.ValueOf(Spider).Type()
					v := reflect.New(t).Elem()

					method := v.MethodByName(request.CallbackMethod)
					this.mmSettings.MMLogger.Printf("%s %d >>> %s\n", runtime.FuncForPC(pc).Name(), this.mmUtils.GetgoID(), request.Url)
					this.engineLog.RequestCount += 1
					response := this.mmRequest.Request(request)
					this.mmSettings.MMLogger.Printf("%s %d <<< %s（Proxy=%s,CurrentUrl=%s）\n", runtime.FuncForPC(pc).Name(), this.mmUtils.GetgoID(), response.Url, response.Request.Proxy, response.CurrentUrl)

					params := []reflect.Value{
						reflect.ValueOf(response),
					}
					method.Call(params)
				}

			}

			//每个任务的间隔时间
			time.Sleep(time.Second * time.Duration(this.mmSettings.Config.ConcurrentRequestSleep))

			if this.closeEngine {
				//关闭引擎
				chum <- 1
				return
			}

		}
	}

}

func (this *mmEngine) _pipelineer(chum chan int) {
	scheduler = MMSchedulerSington()

	//遍历所有的spider，调用每个pipeline的Open函数
	for _, spiderName := range this.spidersName {
		for _, Pipeline := range this.Pipelines {
			Pipeline.Open(spiderName)
		}
	}

	for {

		//获取所有需要被下载器消费的队列，queue+爬虫Name就是队列名
		for _, spiderName := range this.spidersName {

			//拼接queueKey
			queueKey := fmt.Sprintf("item%s", spiderName)
			item := scheduler.ItemPull(queueKey)
			if item != nil {
				//获取到的item下发给所有的pipeline
				for _, Pipeline := range this.Pipelines {
					this.engineLog.PipelineCount += 1
					Pipeline.ProcessItem(item)
				}

			} else {
				//每个任务的间隔时间
				time.Sleep(time.Second * time.Duration(this.mmSettings.Config.ConcurrentPipelineSleep))
			}

			if this.closeEngine {
				//关闭引擎

				//遍历所有的spider，调用每个pipeline的Close函数
				for _, spiderName := range this.spidersName {
					for _, Pipeline := range this.Pipelines {
						Pipeline.Close(spiderName)
					}
				}

				chum <- 1
				return
			}

		}

	}

}
