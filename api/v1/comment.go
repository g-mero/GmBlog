package v1

import (
	"gmeroblog/model"
	"gmeroblog/utils/errmsg"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 新增评论
func AddComment(c *gin.Context) {
	var data model.Comment
	_ = c.ShouldBindJSON(&data)
	userID, bool := c.Get("userID")
	if bool && reflect.TypeOf(userID).String() == "int" {
		data.UserID = userID.(int)
		code := model.AddComment(&data)
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  500,
			"data":    nil,
			"message": "错误：userId获取失败",
		})
	}
}

// 获取评论即回复(由于普通用户都可访问所以建议ps做限制)
func GetComments(c *gin.Context) {
	artId, _ := strconv.Atoi(c.Query("article_id"))
	pn, _ := strconv.Atoi(c.Query("pagenum"))
	var data []model.CmtWithUserWithReplys
	var total int64
	code := errmsg.ERROR
	if pn <= 0 {
		pn = 1
	}
	if artId > 0 {
		data, total, code = model.GetCommentsWithUserByArtID(8, pn, artId)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

func GetReplys(c *gin.Context) {
	cmtId, _ := strconv.Atoi(c.Query("comment_id"))
	pn, _ := strconv.Atoi(c.Query("pagenum"))
	// 限制pagesize为8
	var ps = 8
	var data []model.CmtWithUserWithReplys
	code := errmsg.ERROR
	if pn <= 0 {
		pn = 1
	}
	if cmtId > 0 {
		data, code = model.GetReplysWithUserByCommentID(ps, pn, cmtId)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 获取评论即回复(未登录)
func GetCommentsLimit(c *gin.Context) {
	artId, _ := strconv.Atoi(c.Query("article_id"))
	var data []model.CmtWithUserWithReplys
	var total int64
	code := errmsg.ERROR
	if artId > 0 {
		data, total, code = model.GetCommentsByArtID(6, 1, artId)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}
