package memory

/**
下载器中间件
 */
type MMDownloaderMiddleware interface {

	/**
	参数:
	request (MemoryRequest 对象) – 处理的request
	spider  (MMSpider 对象) – 该request对应的spider

	 */
	ProcessRequest(request MemoryRequest,spider MMSpider) MemoryRequest  //当每个request通过下载中间件时，该方法被调用


	/**
	参数:
	request  (MemoryRequest 对象) – response所对应的request
	response (MemoryResponse 对象) – 被处理的response
	spider   (MMSpider 对象) – response所对应的spider
	 */
	ProcessResponse(request MemoryRequest,response MemoryResponse,spider MMSpider) MemoryResponse //当每个response通过下载中间件时，该方法被调用



}
