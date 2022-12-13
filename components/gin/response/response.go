// Package response Package response 提供统一的 JSON 返回结构，可以通过配置设置具体返回的 code 字段为 int 或者 string
package response

import (
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 响应体
type Response struct {
	Code int         `json:"code"`    // 响应码
	Msg  string      `json:"message"` // 响应消息
	Data interface{} `json:"data"`    // 响应数据
}

// OK 返回 HTTP 状态码为 200 的统一成功结构
func OK(c *gin.Context, data interface{}) {
	Respond(c, http.StatusOK, "", data, nil)
}

// OKMsg 带msg成功消息
func OKMsg(c *gin.Context, message string, data interface{}) {
	Respond(c, http.StatusOK, message, data, nil)
}

// Err 返回 HTTP 状态码为 200 的统一失败结构
func Err(c *gin.Context, errMsg string) {
	Respond(c, http.StatusInternalServerError, errMsg, "", nil)
}

// ErrMsg 返回 HTTP 状态码为 200 的统一失败结构
func ErrMsg(c *gin.Context, errMsg string, err error) {
	Respond(c, http.StatusInternalServerError, errMsg, "", err)
}

// Respond encapsulates c.JSON
// debug mode respond indented json
func Respond(c *gin.Context, code int, msg string, data interface{}, err error) {
	resp := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.JSON(code, resp)
}

func JsonBindError(c *gin.Context, err error) {
	if err.Error() == consts.EOFError {
		Err(c, consts.EOFErrorMsg)
	} else {
		Err(c, err.Error())
	}
}
