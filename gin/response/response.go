package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate stringer -type=StatusCode
type StatusCode int

const (
	Success StatusCode = 0
	Failure StatusCode = 1
)

var (
	emptyData = make(map[string]interface{})
)

// Response gin响应通用结构
type Response struct {
	Code StatusCode  `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Result 返回json格式
func Result(c *gin.Context, code StatusCode, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{code, data, msg})
}

// Ok .
func Ok(c *gin.Context) {
	Result(c, Success, emptyData, Success.String())
}

// OkWithData .
func OkWithData(c *gin.Context, data interface{}) {
	Result(c, Success, data, Success.String())
}

// OkWithMessage .
func OkWithMessage(c *gin.Context, msg string) {
	Result(c, Success, emptyData, msg)
}

// OkWithDetails .
func OkWithDetails(c *gin.Context, data interface{}, msg string) {
	Result(c, Success, data, msg)
}

// Fail .
func Fail(c *gin.Context) {
	Result(c, Failure, emptyData, Failure.String())
}

// FailWithMessage .
func FailWithMessage(c *gin.Context, msg string) {
	Result(c, Failure, emptyData, msg)
}

// FailWithDetails .
func FailWithDetails(c *gin.Context, data interface{}, msg string) {
	Result(c, Failure, data, msg)
}
