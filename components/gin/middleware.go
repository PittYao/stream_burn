package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/helper"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"
	"time"
)

// DefaultGinMiddlewares 默认使用的中间件列表
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		// 记录请求处理日志，最顶层执行
		GinLogger(),
		// 捕获 panic 保存到 context 中由 GinLogger 统一打印， panic 时返回 -1 JSON
		Recovery(),
		// 跨域
		Cors(),
	}

	return m
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// read request body
		var requestBody string
		if method == http.MethodPost || method == http.MethodPut {
			requestBody = ReadRequestBody(c)
		}

		blw := &CustomResponseWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		// 跳过swagger
		if strings.Contains(path, config.C.Swagger.Url) {
			return
		}

		cost := time.Since(start)

		// read response json body
		var responseBody []byte
		if json.Valid(blw.Body.Bytes()) {
			responseBody = json.RawMessage(blw.Body.String())
		} else {
			responseBody = blw.Body.Bytes()
		}

		if method == http.MethodPost || method == http.MethodPut {
			log.L.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("request_body", requestBody),
				zap.Any("response_body", &responseBody),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)

		}

		if method == http.MethodGet || method == http.MethodDelete {
			requestParams := c.Request.URL.RawQuery

			log.L.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("request_params", requestParams),
				zap.Any("response_body", &responseBody),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)
		}

	}
}

// ReadRequestBody 从request中读取body
func ReadRequestBody(c *gin.Context) string {
	var bodyBytes []byte

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Errorf("invalid request body")
	}

	// 新建缓冲区并替换原有Request.body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 当前函数可以使用body内容
	body := bodyBytes
	str := helper.Bytes2str(body)

	// 去除json中的制表符
	compressStr := helper.CompressStr(str)
	return compressStr
}

// Recovery recover掉项目可能出现的panic，并使用zap记录相关日志
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				req, _ := httputil.DumpRequest(c.Request, false)

				println(debug.Stack())
				log.L.Error("[Recovery from panic]",
					zap.Any("error", err),
					zap.Any("request", string(req)),
					zap.ByteString("stack", debug.Stack()),
					zap.Time("runtime", time.Now()),
				)
			}
		}()
		c.Next()
	}
}

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range context.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		if origin != "" {
			context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Ip, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			context.Header("Access-Control-Max-Age", "172800")
			context.Header("Access-Control-Allow-Credentials", "false")
			context.Set("content-type", "application/json")
		}

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "Options Request!")
		}
		//处理请求
		context.Next()
	}
}
