package controller

import (
	"encoding/json"
	"errors"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strings"
)

const (
	Success    = 0
	ParamError = 1
	UserExist  = 2
	SaveError  = 3
)

type Response struct {
	Msg string `json:"msg"`
}

type HTTPError struct {
	Code int
	Err  error
}

type CommonInfo struct {
	PluginAddr string
	PluginPort int
	User       string
	Pwd        string
}

type TokenInfo struct {
	User    string `json:"user" form:"user"`
	Token   string `json:"token" form:"token"`
	Comment string `json:"comment" form:"comment"`
	Status  bool   `json:"status" form:"status"`
}

type TokenResponse struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Count int         `json:"count"`
	Data  []TokenInfo `json:"data"`
}

type OperationResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TokenSearch struct {
	TokenInfo
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type TokenUpdate struct {
	Before TokenInfo `json:"before"`
	After  TokenInfo `json:"after"`
}

type TokenRemove struct {
	Users []TokenInfo `json:"users"`
}

type TokenDisable struct {
	TokenRemove
}

type TokenEnable struct {
	TokenDisable
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

type HandlerFunc func(ctx *gin.Context) (interface{}, error)

func (c *HandleController) MakeHandlerFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		var response plugin.Response
		var err error

		request := plugin.Request{}
		if err := context.BindJSON(&request); err != nil {
			_ = context.Error(&HTTPError{
				Code: http.StatusBadRequest,
				Err:  err,
			})
			return
		}

		jsonStr, err := json.Marshal(request.Content)
		if err != nil {
			_ = context.Error(&HTTPError{
				Code: http.StatusBadRequest,
				Err:  err,
			})
			return
		}

		if request.Op == "Login" {
			content := plugin.LoginContent{}
			err = json.Unmarshal(jsonStr, &content)
			response, err = c.HandleLogin(&content)
		} else if request.Op == "NewProxy" {
			content := plugin.NewProxyContent{}
			err = json.Unmarshal(jsonStr, &content)
			response, err = c.HandleNewProxy(&content)
		} else if request.Op == "Ping" {
			content := plugin.PingContent{}
			err = json.Unmarshal(jsonStr, &content)
			response, err = c.HandlePing(&content)
		} else if request.Op == "NewWorkConn" {
			content := plugin.NewWorkConnContent{}
			err = json.Unmarshal(jsonStr, &content)
			response, err = c.HandleNewWorkConn(&content)
		} else if request.Op == "NewUserConn" {
			content := plugin.NewUserConnContent{}
			err = json.Unmarshal(jsonStr, &content)
			response, err = c.HandleNewUserConn(&content)
		}

		if err != nil {
			log.Printf("handle %s error: %v", context.Request.URL.Path, err)
			var e *HTTPError
			switch {
			case errors.As(err, &e):
				context.JSON(e.Code, &Response{Msg: e.Err.Error()})
			default:
				context.JSON(http.StatusInternalServerError, &Response{Msg: err.Error()})
			}
			return
		} else {
			resStr, _ := json.Marshal(response)
			log.Printf("handle:%v , result: %v", request.Op, string(resStr))
		}

		context.JSON(http.StatusOK, response)
	}
}

func (c *HandleController) MakeManagerFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"UserManage":             ginI18n.MustGetMessage(context, "User Manage"),
			"User":                   ginI18n.MustGetMessage(context, "User"),
			"Token":                  ginI18n.MustGetMessage(context, "Token"),
			"Notes":                  ginI18n.MustGetMessage(context, "Notes"),
			"Search":                 ginI18n.MustGetMessage(context, "Search"),
			"Reset":                  ginI18n.MustGetMessage(context, "Reset"),
			"NewUser":                ginI18n.MustGetMessage(context, "New User"),
			"RemoveUser":             ginI18n.MustGetMessage(context, "Remove User"),
			"DisableUser":            ginI18n.MustGetMessage(context, "Disable User"),
			"EnableUser":             ginI18n.MustGetMessage(context, "Enable User"),
			"Remove":                 ginI18n.MustGetMessage(context, "Remove"),
			"Enable":                 ginI18n.MustGetMessage(context, "Enable"),
			"Disable":                ginI18n.MustGetMessage(context, "Disable"),
			"PleaseInputUserAccount": ginI18n.MustGetMessage(context, "Please Input User Account"),
			"PleaseInputUserToken":   ginI18n.MustGetMessage(context, "Please Input User Token"),
			"PleaseInputUserNotes":   ginI18n.MustGetMessage(context, "Please Input User Notes"),
		})
	}
}

func (c *HandleController) MakeLangFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"User":               ginI18n.MustGetMessage(context, "User"),
			"Token":              ginI18n.MustGetMessage(context, "Token"),
			"Notes":              ginI18n.MustGetMessage(context, "Notes"),
			"Status":             ginI18n.MustGetMessage(context, "Status"),
			"Operation":          ginI18n.MustGetMessage(context, "Operation"),
			"Enable":             ginI18n.MustGetMessage(context, "Enable"),
			"Disable":            ginI18n.MustGetMessage(context, "Disable"),
			"NewUser":            ginI18n.MustGetMessage(context, "New User"),
			"Confirm":            ginI18n.MustGetMessage(context, "Confirm"),
			"Cancel":             ginI18n.MustGetMessage(context, "Cancel"),
			"ConfirmRemoveUser":  ginI18n.MustGetMessage(context, "Confirm to Remove User"),
			"ConfirmDisableUser": ginI18n.MustGetMessage(context, "Confirm to Disable User"),
			"ConfirmEnableUser":  ginI18n.MustGetMessage(context, "Confirm to Enable User"),
			"OperateSuccess":     ginI18n.MustGetMessage(context, "Operate Success"),
			"OperateError":       ginI18n.MustGetMessage(context, "Operate Error"),
			"OperateFailed":      ginI18n.MustGetMessage(context, "Operate Failed"),
			"UserExist":          ginI18n.MustGetMessage(context, "User Exist"),
			"TokenEmpty":         ginI18n.MustGetMessage(context, "Token Can't be Empty"),
			"ShouldCheckUser":    ginI18n.MustGetMessage(context, "Please Check at least One User"),
			"OperationConfirm":   ginI18n.MustGetMessage(context, "Operation Confirm"),
			"EmptyData":          ginI18n.MustGetMessage(context, "Empty Data"),
		})
	}
}

func (c *HandleController) MakeQueryTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {

		search := TokenSearch{}
		search.Limit = 0

		err := context.BindQuery(&search)
		if err != nil {
			return
		}

		var tokenList []TokenInfo
		for _, tokenInfo := range c.Tokens {
			tokenList = append(tokenList, tokenInfo)
		}
		sort.Slice(tokenList, func(i, j int) bool {
			return strings.Compare(tokenList[i].User, tokenList[j].User) < 0
		})

		var filtered []TokenInfo
		for _, tokenInfo := range tokenList {
			if filter(tokenInfo, search.TokenInfo) {
				filtered = append(filtered, tokenInfo)
			}
		}
		if filtered == nil {
			filtered = []TokenInfo{}
		}

		count := len(filtered)
		if search.Limit > 0 {
			start := max((search.Page-1)*search.Limit, 0)
			end := min(search.Page*search.Limit, len(filtered))
			filtered = filtered[start:end]
		}

		context.JSON(http.StatusOK, &TokenResponse{
			Code:  0,
			Msg:   "query Tokens success",
			Count: count,
			Data:  filtered,
		})
	}
}

func filter(main TokenInfo, sub TokenInfo) bool {
	if len(strings.TrimSpace(sub.User)) != 0 {
		if !strings.Contains(main.User, sub.User) {
			return false
		}
	}
	if len(strings.TrimSpace(sub.Token)) != 0 {
		if !strings.Contains(main.Token, sub.Token) {
			return false
		}
	}
	if len(strings.TrimSpace(sub.Comment)) != 0 {
		if !strings.Contains(main.Comment, sub.Comment) {
			return false
		}
	}
	return true
}

func (c *HandleController) MakeAddTokenFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		token := TokenInfo{
			Status: true,
		}
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "add success",
		}
		err := context.BindJSON(&token)
		if err != nil {
			log.Printf("add failed, param error : %v", err)
			response.Success = false
			response.Code = ParamError
			response.Message = "add failed, param error "
			context.JSON(http.StatusOK, &response)
			return
		}
		if _, exist := c.Tokens[token.User]; exist {
			log.Printf("add failed, user %v exist", token.User)
			response.Success = false
			response.Code = UserExist
			response.Message = "add failed, user " + token.User + " exist "
			context.JSON(http.StatusOK, &response)
			return
		}
		token.Comment = ";" + token.Comment
		c.Tokens[token.User] = token

		section, _ := c.IniFile.GetSection("users")
		key, err := section.NewKey(token.User, token.Token)
		key.Comment = token.Comment

		err = c.IniFile.SaveTo(c.ConfigFile)
		if err != nil {
			log.Printf("add failed, error : %v", err)
			response.Success = false
			response.Code = SaveError
			response.Message = "add failed"
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(0, &response)
	}
}

func (c *HandleController) MakeUpdateTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "update success",
		}
		update := TokenUpdate{}
		err := context.BindJSON(&update)
		if err != nil {
			log.Printf("update failed, param error : %v", err)
			response.Success = false
			response.Code = ParamError
			response.Message = "update failed, param error "
			context.JSON(http.StatusOK, &response)
			return
		}

		after := update.After
		before := update.Before

		c.Tokens[after.User] = after

		section, _ := c.IniFile.GetSection("users")
		key, err := section.GetKey(before.User)
		key.Comment = after.Comment
		key.SetValue(after.Token)

		err = c.IniFile.SaveTo(c.ConfigFile)
		if err != nil {
			log.Printf("update failed, error : %v", err)
			response.Success = false
			response.Code = SaveError
			response.Message = "update failed"
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeRemoveTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "remove success",
		}
		remove := TokenRemove{}
		err := context.BindJSON(&remove)
		if err != nil {
			log.Printf("remove failed, param error : %v", err)
			response.Success = false
			response.Code = ParamError
			response.Message = "remove failed, param error "
			context.JSON(http.StatusOK, &response)
			return
		}

		section, _ := c.IniFile.GetSection("users")
		for _, user := range remove.Users {
			delete(c.Tokens, user.User)
			section.DeleteKey(user.User)
		}

		err = c.IniFile.SaveTo(c.ConfigFile)
		if err != nil {
			log.Printf("remove failed, error : %v", err)
			response.Success = false
			response.Code = SaveError
			response.Message = "remove failed"
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeDisableTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "remove success",
		}
		disable := TokenDisable{}
		err := context.BindJSON(&disable)
		if err != nil {
			log.Printf("disable failed, param error : %v", err)
			response.Success = false
			response.Code = ParamError
			response.Message = "disable failed, param error "
			context.JSON(http.StatusOK, &response)
			return
		}

		section, _ := c.IniFile.GetSection("disabled")
		for _, user := range disable.Users {
			section.DeleteKey(user.User)
			token := c.Tokens[user.User]
			token.Status = false
			c.Tokens[user.User] = token
			key, err := section.NewKey(user.User, "disable")
			if err != nil {
				log.Printf("disable failed, error : %v", err)
				response.Success = false
				response.Code = SaveError
				response.Message = "disable failed"
				context.JSON(http.StatusOK, &response)
				return
			}
			key.Comment = "disable user '" + user.User + "'"
		}

		err = c.IniFile.SaveTo(c.ConfigFile)
		if err != nil {
			log.Printf("disable failed, error : %v", err)
			response.Success = false
			response.Code = SaveError
			response.Message = "disable failed"
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeEnableTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "remove success",
		}
		enable := TokenEnable{}
		err := context.BindJSON(&enable)
		if err != nil {
			log.Printf("enable failed, param error : %v", err)
			response.Success = false
			response.Code = ParamError
			response.Message = "enable failed, param error "
			context.JSON(http.StatusOK, &response)
			return
		}

		section, _ := c.IniFile.GetSection("disabled")
		for _, user := range enable.Users {
			section.DeleteKey(user.User)
			token := c.Tokens[user.User]
			token.Status = true
			c.Tokens[user.User] = token
		}

		err = c.IniFile.SaveTo(c.ConfigFile)
		if err != nil {
			log.Printf("enable failed, error : %v", err)
			response.Success = false
			response.Code = SaveError
			response.Message = "enable failed"
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}
