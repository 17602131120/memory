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
        HttpRequest.go
    setting
        settings.go   读取配置
    utils
        mongo.go       mongo工具
        redis.go       redis工具
        util.go        常用工具
    items.go            抓取数据的结构
    models.go           memory.model
    piplines.go         管道
    signal.go           信号
    spider.go           爬虫
    versions.go         版本
    
    

外部依赖
    github.com/PuerkitoBio/goquery
    golang.org/x/net/html
    gopkg.in/mgo.v2                     mongo
    github.com/garyburd/redigo/redis    redis
 -u github.com/go-redis/redis           redis 文档地址 https://godoc.org/github.com/go-redis/redis
    gopkg.in/yaml.v2                    yaml配置文件获取所有属性
    github.com/kylelemons/go-gypsy/yaml yaml配置文件获取某一属性
    第三方log
    github.com/golang/glog
    kafka操作相关
    github.com/Shopify/sarama
    github.com/bsm/sarama-cluster
    

~~~


