package v1

import (
	"gmeroblog/model"
	"gmeroblog/utils/errmsg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// add Diary
func AddDiary(c *gin.Context) {
	var data model.Diary
	_ = c.ShouldBindJSON(&data)
	code := model.CreateDiary(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// get single Diary
func GetSingleDiary(c *gin.Context) {
	time := c.Param("time")
	data, code := model.GetSingleDiary(time)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// search Diarylist
func GetDiary(c *gin.Context) {
	month := c.Query("month")
	year := c.Query("year")

	opt := 0
	value := ""

	if year != "" && month != "" {
		opt = 1
		value = year + "-" + month
	} else if year != "" {
		opt = 2
		value = year
	}

	data, code, total := model.GetDiary(opt, value)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// search DiaryDatelist
func GetDiaryDate(c *gin.Context) {
	month := c.Query("month")
	year := c.Query("year")

	opt := 0
	value := ""

	if year != "" && month != "" {
		opt = 1
		value = year + "-" + month
	} else if year != "" {
		opt = 2
		value = year
	}

	data, code, total := model.GetDiaryDate(opt, value)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// edit Diary
func EditDiary(c *gin.Context) {
	var data model.Diary
	time := c.Param("time")
	c.ShouldBindJSON(&data)
	code := model.EditDiary(time, &data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// delete Diary
func DeleteDiary(c *gin.Context) {
	time := c.Param("time")
	code := model.DeleteDiary(time)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
