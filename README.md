# memory

-   作者：孙国庆
-   邮箱：76251107@qq.com
-   时间：2018/10/17 上午11:20 
-   


#### go 爬虫框架

~~~
基于go 1.11+
    
基于go的爬虫框架 memory

momory
    core
        engine.go     引擎
        scheduler.go  调度
    http
        http.go
    setting
        settings.go   读取配置
    utils
        mongo.go       mongo工具
        redis.go       redis工具
        util.go        常用工具
    items.go            管道数据结构
    models.go           框架提供模型
    middlewares.go      中间件
    piplines.go         管道
    signal.go           信号
    spider.go           爬虫
    versions.go         版本
    
    
附加参数	备  注
-v	显示操作流程的日志及信息，方便检查错误
-u	下载丢失的包，但不会更新已经存在的包
-d	只下载，不安装
-insecure	允许使用不安全的 HTTP 方式进行下载操作


外部依赖
    github.com/PuerkitoBio/goquery
    golang.org/x/net/html
    gopkg.in/mgo.v2                     mongo
    github.com/garyburd/redigo/redis    redis
    github.com/go-redis/redis           redis 文档地址 https://godoc.org/github.com/go-redis/redis
    gopkg.in/yaml.v2                    yaml配置文件获取所有属性
    github.com/kylelemons/go-gypsy/yaml yaml配置文件获取某一属性
    
    github.com/golang/glog 第三方log
    github.com/Shopify/sarama
    github.com/bsm/sarama-cluster
    

~~~


