package model

import (
	"encoding/json"
	"gmeroblog/utils/errmsg"
	"gmeroblog/utils/static"

	"gorm.io/gorm"
)

type Settings struct {
	ID      int    `gorm:"primarykey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Content string `gorm:"not null" json:"content"`
	Role    string `gorm:"not null" json:"role"`
}

type OutSettings struct {
	ID      int               `json:"id"`
	Name    string            `json:"name"`
	Content map[string]string `json:"content"`
}

var initSettings = []Settings{
	{ID: 1, Name: "base_settings", Content: jsonString(map[string]string{
		"site_name": "Gmero's Blog",
		"site_url":  "https://www.gmero.com",
		"site_logo": "/assets/img/logo.svg",
	}), Role: "base"},
	{ID: 2, Name: "seo_settings", Content: jsonString(map[string]string{
		"seo_desc":     "a blog",
		"seo_keywords": "blog,some",
	}), Role: "base"},
	{ID: 3, Name: "admin_settings", Content: jsonString(map[string]string{
		"admin_path":                 "admin",
		"admin_github_id":            "",
		"admin_github_client_id":     "",
		"admin_github_client_secret": "",
	}), Role: "base"},
	{ID: 4, Name: "mail_settings", Content: jsonString(map[string]string{
		"mail_host":     "",
		"mail_port":     "25",
		"mail_username": "",
		"mail_password": "",
	}), Role: "base"},
	{ID: 10, Name: "stay_for_base", Content: "", Role: "base_stay"},
}

func jsonString(arr map[string]string) string {
	bytes, _ := json.Marshal(arr)
	return string(bytes)
}

func saveSets(sets ...OutSettings) {
	for _, set := range sets {
		updateSet(set.Content)
	}
}

func updateSet(content map[string]string) {
	for k, v := range content {
		static.Set("BASE_"+k, v)
	}
}

// 初始化，读取数据库并检查设置项是否完整，最后写入内存
func InitSet() int {
	var setInit = func(set Settings) int {
		var tmp Settings
		err := db.Where("id = ?", set.ID).First(&tmp).Error
		if err == gorm.ErrRecordNotFound {
			err := db.Create(&set).Error
			if err != nil {
				return errmsg.ERROR
			}
		}
		// 修复create失败的问题
		if tmp.Name != set.Name {
			err := db.Save(&set).Error
			if err != nil {
				return errmsg.ERROR
			}
		}
		return errmsg.SUCCES
	}

	for _, v := range initSettings {
		if setInit(v) == errmsg.ERROR {
			return errmsg.ERROR
		}
	}
	// 写入内存
	saveSets(getsets("base")...)
	return errmsg.SUCCES
}

// 获取设置
func GetSettings(role string) ([]OutSettings, int) {
	var settings []Settings
	err := db.Where("role = ?", role).Find(&settings).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	var outset []OutSettings
	for _, value := range settings {
		var tmp OutSettings
		var data map[string]string
		if err := json.Unmarshal([]byte(value.Content), &data); err != nil {
			return nil, errmsg.ERROR
		}
		tmp.Name = value.Name
		tmp.Content = data
		tmp.ID = value.ID
		outset = append(outset, tmp)
	}
	return outset, errmsg.SUCCES
}

func getsets(role string) []OutSettings {
	sets, _ := GetSettings(role)
	return sets
}

// 修改设置
func EditSettings(data OutSettings) int {
	var set Settings
	if err := db.Where("id = ?", data.ID).First(&set).Error; err != nil || set.Name != data.Name {
		return errmsg.ERROR
	}
	set.Content = jsonString(data.Content)
	err := db.Save(&set).Error
	if err != nil {
		return errmsg.ERROR
	}
	updateSet(data.Content)
	return errmsg.SUCCES
}
