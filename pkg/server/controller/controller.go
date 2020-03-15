package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg string `json:"msg"`
}

type HTTPError struct {
	Code int
	Err  error
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

type HandlerFunc func(ctx *gin.Context) (interface{}, error)

func MakeGinHandlerFunc(handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			log.Printf("handle %s error: %v", ctx.Request.URL.Path, err)
			switch e := err.(type) {
			case *HTTPError:
				ctx.JSON(e.Code, &Response{Msg: e.Err.Error()})
			default:
				ctx.JSON(500, &Response{Msg: err.Error()})
			}
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}
