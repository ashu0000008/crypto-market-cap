package main

import (
	"github.com/ashu0000008/crypto-market-cap/account"
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
	engine_https.Use(HeaderHandler())
	engine_https.Any("/", webRoot)
	engine_https.GET("/info/list", getCryptoList)
	engine_https.GET("/crypto/rank", getCryptoRank)
	engine_https.GET("/percent/:symbol", getCryptoPercent)
	engine_https.GET("/percent/:symbol/history", getCryptoPercentHistory)
	engine_https.GET("/platforms/summary", getCryptoPlatformsSummary)
	engine_https.GET("/platform/:platform", getCryptoPlatformInfo)

	engine_https.GET("/favorite", getFavorite)
	engine_https.POST("/favorite", postFavorite)
	engine_https.DELETE("/favorite", deleteFavorite)

	go engine_https.RunTLS(":443", path+"cert.pem", path+"privkey.pem")
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

func HeaderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceId := c.Request.Header.Get("deviceId")

		print("--" + c.Request.Method + c.Request.RequestURI + "-------")
		print("deviceId---------------:" + deviceId + "\r\n")
		account.Check2AddUser(deviceId)

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
	page := context.GetInt("page")
	size := context.GetInt("size")
	context.String(http.StatusOK, impl.GetCryptoListImpl(page, size))
}

func getCryptoRank(context *gin.Context) {
	pageString := context.Query("page")
	sizeString := context.Query("size")

	var page int
	var size int
	var err error

	if "" == pageString {
		page = 0
	} else {
		page, err = strconv.Atoi(pageString)
		if err != nil {
			page = 0
		}
	}
	if "" == sizeString {
		size = 0
	} else {
		size, err = strconv.Atoi(sizeString)
		if err != nil {
			size = 0
		}
	}

	if page == 0 && size == 0 {
		page = 0
		size = 1000
	}

	context.String(http.StatusOK, impl.GetCryptoRankImpl(page, size))
}

func getCryptoPercent(context *gin.Context) {
	symbol := context.Param("symbol")
	percent := impl.GetPercent(symbol)
	context.String(http.StatusOK, strconv.FormatFloat(percent, 'f', -1, 64))
}

func getCryptoPercentHistory(context *gin.Context) {
	symbol := context.Param("symbol")
	context.String(http.StatusOK, impl.GetPercentHistory(symbol))
}

func getCryptoPlatformsSummary(context *gin.Context) {
	context.String(http.StatusOK, impl.GetCryptoPlatformsSummaryImpl())
}

func getCryptoPlatformInfo(context *gin.Context) {
	platform := context.Param("platform")
	context.String(http.StatusOK, impl.GetPlatformInfo(platform))
}

func getFavorite(context *gin.Context) {
	device := context.GetHeader("deviceId")
	context.String(http.StatusOK, account.ApiGetFavorite(device))
}

func postFavorite(context *gin.Context) {
	device := context.GetHeader("deviceId")
	symbol := context.Param("symbol")
	success := account.ApiAddFavorite(device, symbol)
	if success {
		context.String(http.StatusOK, "")
	} else {
		context.String(http.StatusInternalServerError, "")
	}
}

func deleteFavorite(context *gin.Context) {
	device := context.GetHeader("deviceId")
	symbol := context.Param("symbol")
	success := account.ApiDeleteFavorite(device, symbol)
	if success {
		context.String(http.StatusOK, "")
	} else {
		context.String(http.StatusInternalServerError, "")
	}
}
