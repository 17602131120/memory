package http

import (
	"fmt"
	"github.com/sunguoguo/memory"
	"github.com/sunguoguo/memory/setting"
	"github.com/sunguoguo/memory/utils"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"time"
)

type HttpRequest struct{}

//网络请求
func (this *HttpRequest) Request(request memory.MemoryRequest) *memory.MemoryResponse {

	//time.Sleep(300 * time.Millisecond)
	//time.Sleep(1 * time.Second)

	response := new(memory.MemoryResponse)
	response.KeyItemQueue = fmt.Sprintf("item%s", request.CallbackStuct)
	response.KeyRequestQueue = fmt.Sprintf("queue%s", request.CallbackStuct)
	response.Url = request.Url
	response.State = false

	mRedis := new(utils.MMRedis)
	settings := setting.MMSettingsSington()

	mmUserAgent := mRedis.UserAgentSrandmember(settings.Config.UserAgentkey, false)
	mmCookie := mRedis.CookieSrandmember(settings.Config.Cookiekey)
	mmProxy := mRedis.ProxySrandmember(settings.Config.Proxykey)
	request.Proxy = mmProxy.Proxy
	response.Request = request
	var client *http.Client

	if mmProxy.NonProxy {
		client = &http.Client{}
	} else {
		proxyUrl, _ := url.Parse(fmt.Sprintf("http://%s:%s", mmProxy.Ip, mmProxy.Port))
		client = &http.Client{Transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxyUrl),
			ResponseHeaderTimeout: time.Second * 30,
		}}
	}

	req, err := http.NewRequest("GET", request.Url, nil)
	if mmProxy.NonProxy == false && mmProxy.Auth {
		req.SetBasicAuth(mmProxy.Username, mmProxy.Password)
	}
	//req.SetBasicAuth("786251107","oq1fdb7w")
	//req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br") //gzip
	//req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	//req.Header.Add("Cache-Control", "no-cache")
	//req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Host", "")
	req.Header.Add("Cookie", mmCookie.Val)
	//req.Header.Add("Pragma", "no-cache")
	//req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", mmUserAgent.Val)

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
