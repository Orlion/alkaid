package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	code2StatusCodeMap map[int]int
}

func NewResponse(code2StatusCodeMap map[int]int) *Response {
	return &Response{
		code2StatusCodeMap:code2StatusCodeMap,
	}
}

func (r *Response)code2StatusCode(code int) int {
	if statusCode, exist := r.code2StatusCodeMap[code]; exist {
		return statusCode
	}

	return http.StatusOK
}

func (r *Response)ResultJson(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(r.code2StatusCode(code), gin.H{
		"code": code,
		"msg": msg,
		"data": data,
	})
}
