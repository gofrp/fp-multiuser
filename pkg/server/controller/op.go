package controller

import (
	"net/http"

	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
)

type OpController struct {
	tokens map[string]string
}

func NewOpController(tokens map[string]string) *OpController {
	return &OpController{
		tokens: tokens,
	}
}

func (c *OpController) Register(engine *gin.Engine) {
	engine.POST("/handler", MakeGinHandlerFunc(c.HandleLogin))
}

func (c *OpController) HandleLogin(ctx *gin.Context) (interface{}, error) {
	var r plugin.Request
	var content plugin.LoginContent
	r.Content = &content
	if err := ctx.BindJSON(&r); err != nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}

	var res plugin.Response
	token := content.Metas["token"]
	if len(c.tokens) == 0 {
		res.Unchange = true
	} else if content.User == "" || token == "" {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
	} else if c.tokens[content.User] == token {
		res.Unchange = true
	} else {
		res.Reject = true
		res.RejectReason = "invalid meta token"
	}
	return res, nil
}
