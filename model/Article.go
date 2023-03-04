package model

import (
	"fmt"
	"gmeroblog/utils/errmsg"
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID          int            `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Title       string         `gorm:"not null" json:"title"`
	Cid         int            `gorm:"not null;index" json:"cid"`
	Desc        string         `json:"desc"`
	Content     string         `gorm:"not null" json:"content"`
	HTMLContent string         `gorm:"not null" json:"html_content"`
	Img         string         `json:"img"`
	TextCount   int            `gorm:"not null" json:"text_count"`
	Sees        int            `gorm:"not null" json:"sees"`
	Extra       string         `json:"extra"`
}

type Article_with_cate struct {
	Article
	Cate cate
}

type cate struct {
	Name    string `json:"cate_name"`
	Slug    string `json:"cate_slug"`
	GroupId int    `json:"groupId"`
}

// create article
func CreateArt(data *Article) int {
	if data.Title == "" || data.Content == "" || data.Cid <= 0 {
		return errmsg.ERROR
	}

	data.TextCount = len([]rune(data.Content))
	data.Sees = 0
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// 获取单个文章，这里会返回文章信息以及对应分类信息
func GetSingleArt(id int) (Article_with_cate, int) {
	var art Article_with_cate
	var tmp Article
	err = db.Where("id = ?", id).First(&tmp).Error
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	tmp.Sees++
	// 这里使用UpdateColumn防止updateat时间被更新
	db.Model(&tmp).UpdateColumn("sees", tmp.Sees)
	art.Article = tmp
	var cate cate
	var tmpc Category
	err = db.Where("id = ?", art.Cid).First(&tmpc).Error
	if err != nil {
		return art, errmsg.ERROR
	}
	cate.Name = tmpc.Name
	cate.Slug = tmpc.Slug
	cate.GroupId = tmpc.GroupId
	art.Cate = cate
	return art, errmsg.SUCCES
}

// get Article list
func GetArt(pageSize int, pageNum int, cid int) ([]Article, int, int64) {
	var artList []Article
	var total int64
	if cid <= 0 {
		err = db.Model(&artList).Order("id DESC").Count(&total).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&artList).Error
	} else {
		err = db.Model(&artList).Where("cid = ?", cid).Count(&total).Order("id DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&artList).Error
	}
	if err != nil {
		fmt.Println(err)
		return nil, errmsg.ERROR, 0
	}
	return artList, errmsg.SUCCES, total
}

// edit Article
func EditArt(id int, data *Article) int {
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["html_content"] = data.HTMLContent
	maps["img"] = data.Img
	maps["text_count"] = len([]rune(data.Content))
	err = db.Model(&Article{}).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// update HtmlContent
func UpdateHtmlContent(id int, data *Article) int {
	var maps = make(map[string]interface{})
	maps["desc"] = data.Desc
	maps["html_content"] = data.HTMLContent
	// 防止updateat被更新
	err = db.Model(&Article{}).Where("id = ?", id).UpdateColumns(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// delete Article
func DeleteArt(id int) int {
	var art Article
	err = db.Where("id = ?", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}
