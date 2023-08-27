package controller

import (
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

type HandleController struct {
	CommonInfo CommonInfo
	Tokens     map[string]TokenInfo
	ConfigFile string
	IniFile    *ini.File
}

func NewHandleController(config *HandleController) *HandleController {
	return config
}

func (c *HandleController) Register(engine *gin.Engine) {
	engine.Delims("${", "}")
	engine.LoadHTMLGlob("./assets/templates/*")
	engine.POST("/handler", c.MakeHandlerFunc())

	var group *gin.RouterGroup
	if len(c.CommonInfo.User) != 0 {
		group = engine.Group("/", gin.BasicAuthForRealm(gin.Accounts{
			c.CommonInfo.User: c.CommonInfo.Pwd,
		}, "Restricted"))
	} else {
		group = engine.Group("/")
	}
	group.Static("/static", "./assets/static")
	group.GET("/", c.MakeManagerFunc())
	group.GET("/lang", c.MakeLangFunc())
	group.GET("/tokens", c.MakeQueryTokensFunc())
	group.POST("/add", c.MakeAddTokenFunc())
	group.POST("/update", c.MakeUpdateTokensFunc())
	group.POST("/remove", c.MakeRemoveTokensFunc())
	group.POST("/disable", c.MakeDisableTokensFunc())
	group.POST("/enable", c.MakeEnableTokensFunc())
}

func (c *HandleController) HandleLogin(content *plugin.LoginContent) (plugin.Response, error) {
	token := content.Metas["token"]
	user := content.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) HandleNewProxy(content *plugin.NewProxyContent) (plugin.Response, error) {
	token := content.User.Metas["token"]
	user := content.User.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) HandlePing(content *plugin.PingContent) (plugin.Response, error) {
	token := content.User.Metas["token"]
	user := content.User.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) HandleNewWorkConn(content *plugin.NewWorkConnContent) (plugin.Response, error) {
	token := content.User.Metas["token"]
	user := content.User.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) HandleNewUserConn(content *plugin.NewUserConnContent) (plugin.Response, error) {
	token := content.User.Metas["token"]
	user := content.User.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) JudgeToken(user string, token string) (plugin.Response, error) {
	var res plugin.Response
	if len(c.Tokens) == 0 {
		res.Unchange = true
	} else if user == "" || token == "" {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
	} else if info, exist := c.Tokens[user]; exist {
		if !info.Status {
			res.Reject = true
			res.RejectReason = "user " + user + " is disabled"
		} else {
			if info.Token != token {
				res.Reject = true
				res.RejectReason = "invalid meta token for user " + user + ""
			} else {
				res.Unchange = true
			}
		}
	} else {
		res.Reject = true
		res.RejectReason = "user " + user + " not exist"
	}
	return res, nil
}
