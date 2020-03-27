package http

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
)

type Response struct {
	code2StatusCodeMap map[int]int
	viewSet            *jet.Set
}

func NewResponse(code2StatusCodeMap map[int]int, dirs ...string) *Response {
	return &Response{
		code2StatusCodeMap: code2StatusCodeMap,
		viewSet:            jet.NewHTMLSet(dirs...),
	}
}

func (r *Response) code2StatusCode(code int) int {
	if statusCode, exist := r.code2StatusCodeMap[code]; exist {
		return statusCode
	}

	return http.StatusOK
}

func (r *Response) ResultJson(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(r.code2StatusCode(code), gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func (r *Response) View(ctx *gin.Context, code int, viewName string, vars jet.VarMap) (err error) {
	t, err := r.viewSet.GetTemplate(viewName)
	if err != nil {
		err = fmt.Errorf("视图获取失败: [%w]", err)
		return
	}

	var w bytes.Buffer
	if err = t.Execute(&w, vars, nil); err != nil {
		err = fmt.Errorf("视图渲染失败: [%w]", err)
		return
	}

	ctx.Data(r.code2StatusCode(code), "text/html", w.Bytes())

	return
}
