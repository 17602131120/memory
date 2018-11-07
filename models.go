package memory

import (
	"golang.org/x/net/html"
	"time"
)

type MemoryProxy struct {
	//localhost 本地请求，不设置代理
	//ip:port,username:password
	Proxy    string // 代理原始字符串
	Ip       string //ip
	Port     string //port
	NonProxy bool   //本次是否需要设置代理
	Auth     bool   //是否需要账号密码
	Username string
	Password string
}
type MemoryEngineLog struct {
	StartTime     time.Time //引擎开始时间
	FinishTime    time.Time //引擎结束时间
	MainID        int       //线程id
	FinishReason  string    //结束原因
	RequestCount  int       //请求总次数
	PipelineCount int       //写入管道总数量

}
type MemoryRequest struct {
	Url        string            //请求URL
	Proxy      string            //本次请求使用的代理
	Headers    map[string]string //请求头
	Meta       map[string]string //自定义meta
	Depth      int               //获取到该任务的网页深度，首次请求为1
	RequestNum int               //请求尝试次数，首次请求为1
	//Callback   func(response MemoryResponse)
	CallbackStuct  string //callback 结构体 解析函数的字符串，由于go无法实现字符串对函数的序列化反序列化执行
	CallbackMethod string //callback 方法

	Encoding   string
	Priority   int  //优先级
	DontFilter bool //是否不过滤，默认false 过滤

}
type MemoryResponse struct {
	State           bool       //自定义成功失败状态
	StatusCode      int        //响应状态码
	RetryNum        int        //本次请求内部请求的次数
	HtmlNode        *html.Node //htmlNode
	Msg             string     //自定义响应日志
	Url             string     //请求URL
	CurrentUrl      string     //请求响应时URL
	Request         MemoryRequest
	KeyRequestQueue string //
	KeyItemQueue    string //
}

type MemoryCookie struct {
	Auth     bool   //是否登录
	Username string //账号名称
	Password string //密码
	Val      string //Cookie的具体值
	Remark   string //备注

}

type MemoryUserAgent struct {
	Mobile   bool   //是否移动端
	Platform string //平台
	Browser  string //浏览器类型
	Val      string //User-Agent的具体值

}
