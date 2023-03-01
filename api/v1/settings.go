package v1

import (
	"gmeroblog/model"
	"gmeroblog/utils/errmsg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// search Settings
func GetSettings(c *gin.Context) {
	role := c.Query("role")

	data, code := model.GetSettings(role)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// edit settings
func EditSettings(c *gin.Context) {
	var data model.OutSettings
	c.ShouldBindJSON(&data)
	code := model.EditSettings(data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
