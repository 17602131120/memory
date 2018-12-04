package app

import "github.com/sunguoguo/memory"

type CustomerDownloaderMiddleWare struct {
	memory.MMDownloaderMiddleware

}

func (this *CustomerDownloaderMiddleWare) ProcessRequest(request memory.MemoryRequest,spider memory.MMSpider) memory.MemoryRequest  {

	request.Proxy="localhost"


	//req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br") //gzip
	//req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	//req.Header.Add("Cache-Control", "no-cache")
	//req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Host", "")
	//req.Header.Add("Cookie", mmCookie.Val)
	//req.Header.Add("Pragma", "no-cache")
	//req.Header.Add("Upgrade-Insecure-Requests", "1")
	//req.Header.Add("User-Agent", mmUserAgent.Val)

	request.Headers["User-Agent"]="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36"
	request.Headers["Cookie"]="_lxsdk_cuid=1648e9eecdac8-0b185e0fc7e2e9-163b6953-fa000-1648e9eecdac8; _lxsdk=1648e9eecdac8-0b185e0fc7e2e9-163b6953-fa000-1648e9eecdac8; _hc.v=1ed9ddec-cf41-316a-8ec6-2e0b23da35a4.1531401137; s_ViewType=10; _dp.ac.v=1256daa7-6107-4567-9b7f-ca2fb4814882; ua=%E6%96%AD%E6%A1%A5%E6%AE%8B%E9%9B%AA_9287; ctu=89333a78b2aa32adb02ed95af63265e2140105249277d131e3294dfa2667f15c; switchcityflashtoast=1; _tr.u=28PoV8Ds7MrorNl4; looyu_id=b5285a366b95cfe82ace4b685ae8418038_51868%3A1; aburl=1; __utma=1.1616048375.1537188932.1537188932.1537188932.1; __utmz=1.1537188932.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); cityid=2; cy=1; cye=shanghai; Hm_lvt_dbeeb675516927da776beeb1d9802bd4=1538056542,1538123829,1538813206,1540632140; _adwp=169583271.3937982438.1537503412.1537503412.1540632140.2; m_flash2=1; source=m_browser_test_33; default_ab=shop%3AA%3A1%7Cindex%3AA%3A1%7CshopList%3AA%3A1%7Cmap%3AA%3A1; pvhistory=6L+U5ZuePjo8L2Vycm9yL2Vycm9yX3BhZ2U+OjwxNTQzNTY1NzYyNzg4XV9b; _lxsdk_s=16779502312-caf-f9-52d%7C%7C8"



	return request


}

func (this *CustomerDownloaderMiddleWare) ProcessResponse(request memory.MemoryRequest,response memory.MemoryResponse,spider memory.MMSpider) memory.MemoryResponse {



	return response
}

