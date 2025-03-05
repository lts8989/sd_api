package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
)

var Trans ut.Translator

type Context struct {
	*gin.Context
}
type Response struct {
	Code        RespCode `json:"code"`
	MessageData string   `json:"message_data"`
	Data        any      `json:"data"`
}

type RespCode int

const (
	RespCodeSuccess    RespCode = 0
	RespCodeError      RespCode = 1
	RespCodeParamError RespCode = 2
	RespCodeBizError   RespCode = 500
)

func (ctx *Context) Success(data any, msg ...string) {
	message := ""
	if len(msg) > 0 {
		for _, m := range msg {
			message = fmt.Sprint("%s %s", message, m)
		}
	}
	ctx.JSONP(http.StatusOK, &Response{
		Code:        RespCodeSuccess,
		MessageData: message,
		Data:        data,
	})
}

func (ctx *Context) Error(err error) {
	var errs validator.ValidationErrors
	switch {
	case errors.As(err, &errs):
		ctx.ParamError(errs)
		return
	default:
	}

	ctx.JSONP(http.StatusOK, &Response{
		Code:        RespCodeError,
		MessageData: err.Error(),
		Data:        nil,
	})
}

func (ctx Context) ParamError(errs validator.ValidationErrors) {
	errMsg := make(map[string]string)
	for _, e := range errs {
		errMsg[e.Field()] = e.Translate(Trans)

	}

	ctx.JSONP(http.StatusOK, &Response{
		Code:        RespCodeParamError,
		MessageData: "",
		Data:        errMsg,
	})
}

type HandlerFunc = func(ctx *Context)

func Build(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		//这里可以加一下参数统一校验的逻辑
		ctx := &Context{
			Context: c,
		}
		h(ctx)
	}
}
