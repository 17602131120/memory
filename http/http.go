package http

import (
	"fmt"
	"github.com/sunguoguo/memory"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"time"
)

type MMHttp struct{}

//网络请求
func (this *MMHttp) Request(request *memory.MemoryRequest, ) *memory.MemoryResponse {

	//time.Sleep(300 * time.Millisecond)
	//time.Sleep(1 * time.Second)

	//构建response
	response := new(memory.MemoryResponse)

	response.KeyItemQueue = fmt.Sprintf("item%s", request.CallbackStuct)
	response.KeyRequestQueue = fmt.Sprintf("queue%s", request.CallbackStuct)
	response.Url = request.Url
	response.State = false
	response.Request = *request

	var client *http.Client
	var TargetUrl string

	//获取请求类中的代理，并进行实例化
	mmProxy :=new(memory.MemoryProxy)
	mmProxy.SetAttrs(request.Proxy)

	if mmProxy.ProxyType==10{
		//代理转发请求
		client = &http.Client{}
		//拼接代理转发服务的URL
		TargetUrl=fmt.Sprintf("http://%s:%s/?%s",mmProxy.Ip,mmProxy.Port,request.Url)

	}else{
		//代理或本地请求
		TargetUrl= request.Url


		if mmProxy.ProxyType==0 {
			client = &http.Client{}
		} else {
			//代理请求
			proxyUrl, _ := url.Parse(fmt.Sprintf("http://%s:%s", mmProxy.Ip, mmProxy.Port))
			client = &http.Client{Transport: &http.Transport{
				Proxy:                 http.ProxyURL(proxyUrl),
				ResponseHeaderTimeout: time.Second * 30,
			}}


		}

	}

	//发起请求
	req, err := http.NewRequest("GET", TargetUrl, nil)

	if mmProxy.ProxyType == 6 {
		req.SetBasicAuth(mmProxy.Username, mmProxy.Password)
	}
	//req.SetBasicAuth("786251107","oq1fdb7w")
	for k,v := range request.Headers{
		req.Header.Add(k, v)
	}

	if err != nil {
		//失败
		response.Msg = fmt.Sprintf("http.NewRequest is error:%s", err)
	} else {

		res, err := client.Do(req)
		if err == nil {
			defer res.Body.Close()
			//成功
			//body, err := ioutil.ReadAll(res.Body)
			//if err != nil {

			//}
			// []int{403, 405, 500, 502, 503, 504, 408}

			htmlNode, err := html.Parse(res.Body)
			if err != nil {
				//转换时失败
				response.Msg = fmt.Sprintf("html.Parse is error:%s", err)

			} else {

				//log.Println(res.Request.URL.String())

				response.StatusCode = res.StatusCode
				response.CurrentUrl = res.Request.URL.String()
				response.HtmlNode = htmlNode
				response.State = true

			}
		} else {

			response.Msg = fmt.Sprintf("client.Do is error:%s", err)
		}

	}





	return response

}
