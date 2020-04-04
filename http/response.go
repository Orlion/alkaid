package http

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
)

type Seo struct {
	Title       string
	Description string
	Keywords    string
	Author      string
}

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
		"message":  msg,
		"data": data,
	})
}

func (r *Response) View(ctx *gin.Context, code int, viewName string, vars jet.VarMap, seo *Seo) {
	t, err := r.viewSet.GetTemplate(viewName)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("视图加载失败"))
		return
	}

	statusCode := r.code2StatusCode(code)
	vars.Set("statusCode", statusCode)
	vars.Set("seo", seo)
	vars.Set("requestUrl", ctx.Request.URL)

	var w bytes.Buffer
	if err = t.Execute(&w, vars, nil); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("视图渲染失败"))
		return
	}

	ctx.Data(statusCode, "text/html", w.Bytes())
	return
}
