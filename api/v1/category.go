package v1

import (
	"gmeroblog/model"
	"gmeroblog/utils/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// add Category
func AddCategory(c *gin.Context) {
	var data model.Category
	_ = c.ShouldBindJSON(&data)
	code := model.CheckCategory(data)
	if code == errmsg.SUCCES {
		code = model.CreateCate(&data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// get single Category
func GetSingleCate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetSingleCate(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"massage": errmsg.GetErrMsg(code),
	})
}

// search Category
func GetCate(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	if pageSize == 0 {
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = 1
	}
	data := model.GetCate(pageSize, pageNum)
	code := errmsg.SUCCES
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"massage": errmsg.GetErrMsg(code),
	})
}

// search Category tree
func GetCateTree(c *gin.Context) {
	data := model.GetCates()
	code := errmsg.SUCCES
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"massage": errmsg.GetErrMsg(code),
	})
}

// edit Category
func EditCate(c *gin.Context) {
	var data model.Category
	id, _ := strconv.Atoi(c.Param("id"))
	data.ID = id
	c.ShouldBindJSON(&data)
	code := model.CheckCategory(data)
	if code == errmsg.SUCCES {
		model.EditCate(id, &data)
	} else {
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// delete Category
func DeleteCate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteCate(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
