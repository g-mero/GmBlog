package v1

import (
	"gmeroblog/middleware"
	"gmeroblog/model"
	"gmeroblog/utils/errmsg"
	"gmeroblog/utils/validater"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// check token
func CheckToken(c *gin.Context) {
	code := errmsg.ERROR
	token := c.Query("token")
	if token == "" {
		code = 500
	} else {
		key, ecode := middleware.CheckToken(token)
		if ecode != errmsg.SUCCES {
			code = 500
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  200,
				"userid":  key.UserID,
				"role":    key.Role,
				"message": errmsg.GetErrMsg(code),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// add user
func AddUser(c *gin.Context) {
	var data model.User
	var msg string
	code := errmsg.ERROR
	_ = c.ShouldBindJSON(&data)
	msg, code = validater.Validate(&data)
	if code != errmsg.SUCCES {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return
	}

	code = model.CheckUser(data.Username)
	if code == errmsg.SUCCES {
		model.CreateUser(&data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		code = errmsg.ERROR_USERNAME_USED
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// search single user
func GetUserInfo(c *gin.Context) {
	username, exist := c.Get("username")
	if exist {
		str, ok := username.(string)
		if ok {
			user, code := model.GetUserByUsername(str)
			var ouser struct {
				Username  string `json:"username"`
				Nickname  string `json:"nickname"`
				AvatarUrl string `json:"avatar_url"`
			}
			ouser.AvatarUrl = user.AvatarUrl
			ouser.Username = user.Username
			ouser.Nickname = user.Nickname
			if code == errmsg.SUCCES {
				c.JSON(http.StatusOK, gin.H{
					"status":  code,
					"data":    ouser,
					"massage": errmsg.GetErrMsg(code),
				})
				return
			}
		}
	}
	c.JSON(500, gin.H{
		"status":  500,
		"massage": errmsg.GetErrMsg(500),
	})
}

// search user list
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	if pageSize == 0 {
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = 1
	}
	data := model.GetUsers(pageSize, pageNum)
	code := errmsg.SUCCES
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"massage": errmsg.GetErrMsg(code),
	})
}

// edit user
func EditUser(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	c.ShouldBindJSON(&data)
	code := model.CheckUser(data.Username)
	if code == errmsg.SUCCES {
		model.EditUser(id, &data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// delete user
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
