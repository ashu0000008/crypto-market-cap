package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/ashu0000008/crypto-market-cap/api/impl"
)

func main() {
	// 初始化引擎
	engine := gin.Default()
	// 注册一个路由和处理函数
	engine.Any("/", webRoot)
	engine.GET("/info/list", getCryptoList)
	engine.GET("/percent/:symbol", getCryptoPercent)
	engine.GET("/platforms/summary", getCryptoPlatformsSummary)
	// 绑定端口，然后启动应用
	engine.Run(":80")
}

/**
* 根请求处理函数
* 所有本次请求相关的方法都在 context 中，完美
* 输出响应 hello, world
 */
func webRoot(context *gin.Context) {
	context.String(http.StatusOK, "hello, world")
}

func getCryptoList(context *gin.Context) {
	context.String(http.StatusOK, impl.GetCryptoListImpl())
}

func getCryptoPercent(context *gin.Context) {
	symbol := context.Param("symbol")
	percent := impl.GetPercent(symbol)
	context.String(http.StatusOK, strconv.FormatFloat(percent, 'f', -1, 64))
}

func getCryptoPlatformsSummary(context *gin.Context) {
	context.String(http.StatusOK, impl.GetCryptoPlatformsSummaryImpl())
}
