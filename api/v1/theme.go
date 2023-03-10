package v1

import (
	"gmeroblog/model"
	"gmeroblog/utils/errmsg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetThemeInfo(c *gin.Context) {

	data, code := model.GetThemeInfo()

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"massage": errmsg.GetErrMsg(code),
	})
}

func GetThemeSettings(c *gin.Context) {
	data, code := model.GetThemeSettings()

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"massage": errmsg.GetErrMsg(code),
	})
}

// edit theme settings
func EditThemeSettings(c *gin.Context) {
	var data model.ThemeSetting
	c.ShouldBindJSON(&data)
	code := model.EditThemeSetting(data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
