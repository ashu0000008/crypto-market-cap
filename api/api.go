package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
	"runtime"
	"strconv"

	"github.com/ashu0000008/crypto-market-cap/api/impl"
)

func main() {

	var host string
	var path string
	if "darwin" != runtime.GOOS {
		host = "https://ashu.xyz"
		path = "/mytls/"
	} else {
		host = "https://localhost"
		path = "/Users/zhouyang/mytls/"
	}

	// 初始化引擎
	engine := gin.Default()
	engine.GET("/*path", func(c *gin.Context) {
		c.Redirect(302, host+c.Param("path"))
	})

	// 绑定端口，然后启动应用
	engine_https := gin.Default()
	engine_https.Use(TlsHandler())
	engine_https.Any("/", webRoot)
	engine_https.GET("/info/list", getCryptoList)
	engine_https.GET("/percent/:symbol", getCryptoPercent)
	engine_https.GET("/platforms/summary", getCryptoPlatformsSummary)

	go engine_https.RunTLS(":443", path+"chain.pem", path+"privkey.pem")
	engine.Run(":80")
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:80",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
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
