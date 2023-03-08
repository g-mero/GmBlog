package v1

import (
	"encoding/json"
	"gmeroblog/middleware"
	"gmeroblog/model"
	"gmeroblog/utils/cache"
	"gmeroblog/utils/errmsg"
	"gmeroblog/utils/static"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var data model.User
	c.ShouldBindJSON(&data)
	var code int
	var token string

	data, code = model.CheckLogin(data.Username, data.Password)

	if code == errmsg.SUCCES {
		token, code = middleware.SetToken(data.ID, 1)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     code,
		"message":    errmsg.GetErrMsg(code),
		"token":      token,
		"nickname":   data.Nickname,
		"avatar_url": data.AvatarUrl,
	})
}

func Logout(c *gin.Context) {
	login_token, _ := c.Get("loginToken")
	code := errmsg.SUCCES
	if login_token != "" {
		cache.Del(login_token.(string))
	} else {
		// 登出失败
		code := errmsg.ERROR
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// github auth
func GithubLogin(c *gin.Context) {
	// 检查是否开启GITHUB认证
	if static.Get("admin_github_client_id") == "" {
		c.HTML(200, "githubAuth", gin.H{
			"title":  "错误",
			"info":   "服务器未开启github认证",
			"status": 500,
		})
		return
	}

	var accessToken = ""

	// 防止重复登录请求githubtoken(这个其实命中率很低，因为登录状态下也不需要重复申请，这里只是防止请求被重复触发)
	login_token, err := c.Cookie("login_token")
	if login_token != "" && err == nil {
		key, code := middleware.CheckToken(login_token)
		if code == errmsg.SUCCES {
			user, code := model.GetUserByUserID(key.UserID)
			if code == errmsg.SUCCES {
				accessToken = user.GithubToken
			}
		}
	}
	if accessToken == "" {
		accessToken = get_github_token(c)
		if accessToken == "" {
			c.HTML(200, "githubAuth", gin.H{
				"title":  "错误",
				"info":   "Github Token 获取失败",
				"status": 500,
			})
			return
		}
	}

	data, code := model.GithubLogin(accessToken)
	if code == 200 {
		token, code := middleware.SetToken(data.ID, data.Role)
		if code == errmsg.SUCCES {
			save_cookies(data, token, 7*24*3600, c)
		} else {
			c.HTML(200, "githubAuth", gin.H{
				"title":  errmsg.GetErrMsg(code),
				"info":   "LoginToken 生成失败",
				"status": code,
			})
			return
		}
		c.HTML(200, "githubAuth", gin.H{
			"title":  errmsg.GetErrMsg(code),
			"info":   errmsg.GetErrMsg(code),
			"status": code,
		})
	} else {
		c.HTML(200, "githubAuth", gin.H{
			"title":  errmsg.GetErrMsg(code),
			"info":   errmsg.GetErrMsg(code),
			"status": code,
		})
	}
}

func get_github_token(c *gin.Context) string {
	gcode := c.Query("code")
	accessTokenURL, _ := url.Parse("https://github.com/login/oauth/access_token")
	values := url.Values{}
	values.Set("client_id", static.Get("admin_github_client_id"))
	values.Set("client_secret", static.Get("admin_github_client_secret"))
	values.Set("code", gcode)
	accessTokenURL.RawQuery = values.Encode()
	req, err := http.NewRequest("POST", accessTokenURL.String(), nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}

	if resp.StatusCode != 200 {
		return ""
	}
	var accessTokenStruct struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&accessTokenStruct); err != nil {
		return ""
	}
	resp.Body.Close()
	return accessTokenStruct.AccessToken
}

// 单位是秒
func save_cookies(user model.User, token string, time int, c *gin.Context) {
	c.SetCookie("login_token", token, time, "/", "", false, false)
	c.SetCookie("nickname", user.Nickname, time, "/", "", false, false)
	c.SetCookie("avatar_url", user.AvatarUrl, time, "/", "", false, false)
}
