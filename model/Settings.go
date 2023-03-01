package model

import (
	"encoding/json"
	"fmt"
	"gmeroblog/utils/errmsg"

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
		"site_name":        "Gmero's Blog",
		"site_url":         "https://www.gmero.com",
		"site_logo":        "/assets/img/logo.svg",
		"site_username":    "GmBlog",
		"site_user_avatar": "https://avatars.githubusercontent.com/u/0",
		"site_user_desc":   "Hello world",
		"site_notice":      "GmBlog欢迎你",
	}), Role: "base"},
	{ID: 2, Name: "seo_settings", Content: jsonString(map[string]string{
		"seo_desc":     "a blog",
		"seo_keywords": "blog,some",
		"seo_footer": `<div class="clearfix px-1"><div class="float-end">
		<a href="https://www.gmero.com" data-bs-toggle="tooltip" data-bs-placement="top"
		title="Gmero'blog" target="_blank"><img data-src="https://img.shields.io/badge/Powered-Gmero-brightgreen" alt="" class="lazy-load" />
		</a>
		</div></div>`,
	}), Role: "base"},
	{ID: 3, Name: "art_settings", Content: jsonString(map[string]string{
		"art_recommend": "",
		"art_top":       "",
	}), Role: "base"},
	{ID: 4, Name: "admin_settings", Content: jsonString(map[string]string{
		"admin_path":                 "admin",
		"admin_github_id":            "",
		"admin_github_client_id":     "",
		"admin_github_client_secret": "",
	}), Role: "base"},
}

var SITE_SETTING map[string]string

func jsonString(arr map[string]string) string {
	bytes, _ := json.Marshal(arr)
	return string(bytes)
}

func mapConnect(sets ...OutSettings) map[string]string {
	out := make(map[string]string)

	for _, set := range sets {
		for k, v := range set.Content {
			out[k] = v
		}
	}
	fmt.Println(out)

	return out
}

func InitSet() int {
	var tmp Settings
	err := db.Where("id = ?", "1").First(&tmp).Error
	if err == gorm.ErrRecordNotFound || tmp.Name != "base_settings" {
		err = db.Create(&initSettings).Error
		if err != nil {
			return errmsg.ERROR
		}
	}

	var setInit = func(set Settings) int {
		var tmp Settings
		err := db.Where("id = ?", set.ID).First(&tmp).Error
		if err == gorm.ErrRecordNotFound || tmp.Name != set.Name {
			err = db.Create(&set).Error
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

	SITE_SETTING = mapConnect(getsets("base")...)
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
	InitSet()
	return errmsg.SUCCES
}
