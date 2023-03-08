package v1

import (
	"gmeroblog/model"
	"gmeroblog/utils/errmsg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 添加一条心情
func AddMood(c *gin.Context) {
	var data model.Mood
	_ = c.ShouldBindJSON(&data)

	code := model.AddMood(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})

}

// 获取一条随机Mood
func GetRandMood(c *gin.Context) {
	data := model.GetMoodRand()

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"message": errmsg.GetErrMsg(200),
	})
}
