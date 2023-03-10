package model

import (
	"encoding/json"
	"gmeroblog/utils/config"
	"gmeroblog/utils/errmsg"
	"gmeroblog/utils/static"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type ThemeInfo struct {
	Name        string `yaml:"name" json:"name"`                 // 主题的唯一标识
	Version     string `yaml:"version" json:"version"`           // 主题版本
	GmVersion   string `yaml:"gm_version" json:"gm_version"`     // 主题适配的博客程序版本
	DisplayName string `yaml:"display_name" json:"display_name"` // 主题显示的名称
	Author      string `yaml:"author" json:"author"`             // 作者名称
	Website     string `yaml:"website" json:"website"`           // 网站
	Repo        string `yaml:"repo" json:"repo"`                 // 仓库
	Desc        string `yaml:"desc" json:"desc"`                 // 主题描述
}

type ThemeSetting struct {
	Name    string   `yaml:"name" json:"name"`
	Label   string   `yaml:"label" json:"label"`
	Content []oneSet `yaml:"content" json:"content"`
}

type oneSet struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
	Label string `yaml:"label" json:"label"`
	Desc  string `yaml:"desc" json:"desc"`
	Type  string `yaml:"type" json:"type"` // text | textarea | radio | select | switch
}

var ThemeName = getInfo().Name
var themePath = "web/theme/" + config.Theme + "/"
var infoPath = themePath + "Theme.yaml"
var setPath = themePath + "Settings.yaml"

// 初始化
func InitTheme() {
	info := getInfo()

	infoByte, err := yaml.Marshal(info)
	if err != nil {
		log.Fatal("[Theme]", err)
	}
	static.Set("theme_info", string(infoByte))

	sets, code := GetThemeSettings()

	if code != errmsg.SUCCES {
		log.Fatal("[Theme]", "获取设置项失败")
	}

	if initThemeDb(sets) != errmsg.SUCCES {
		log.Fatal("[Theme]", "初始化数据库失败")
	}
}

func initThemeDb(sets []ThemeSetting) int {
	var setInit = func(set Settings) int {
		var tmp Settings
		err := db.Where("name = ? AND role = ?", set.Name, set.Role).First(&tmp).Error
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&set).Error
			if err != nil {
				return errmsg.ERROR
			}
		}
		return errmsg.SUCCES
	}

	role := "theme_" + getInfo().Name
	var themeSetsDb []Settings
	for _, v := range sets {
		var tmp Settings
		tmp.Name = v.Name
		tmp.Role = role
		bytes, _ := json.Marshal(v.Content)
		tmp.Content = string(bytes)
		themeSetsDb = append(themeSetsDb, tmp)
		// 写入内存
		saveThemeSetToStatic(v)
	}

	// 数据库处理
	for _, v := range themeSetsDb {
		if setInit(v) != errmsg.SUCCES {
			return errmsg.ERROR
		}
	}

	return errmsg.SUCCES
}

// 从文件获取主题信息
func getInfo() ThemeInfo {
	file, err := os.ReadFile(infoPath)

	if err != nil {
		log.Fatalln("[Theme]", err)
	}

	var info ThemeInfo

	if err := yaml.Unmarshal(file, &info); err != nil {
		log.Fatalln("[Theme]", err)
	}

	return info
}

func GetThemeSettings() ([]ThemeSetting, int) {
	var sets []ThemeSetting

	file, err := os.ReadFile(setPath)
	if err != nil {
		return nil, errmsg.ERROR
	}

	if err := yaml.Unmarshal(file, &sets); err != nil {
		return nil, errmsg.ERROR
	}

	// 将数据库中的数据与设置项合并传值
	var dbSets []Settings
	var newSets []ThemeSetting

	if err := db.Where("role = ?", "theme_"+ThemeName).Find(&dbSets).Error; err != nil {
		return nil, errmsg.ERROR
	}

	for _, v := range dbSets {
		var tmp ThemeSetting
		tmp.Name = v.Name
		var cont []oneSet
		json.Unmarshal([]byte(v.Content), &cont)
		tmp.Content = cont
		newSets = append(newSets, tmp)
	}

	var assign = func(old []oneSet, new []oneSet) {
		for k, v := range old {
			for _, ve := range new {
				if v.Name == ve.Name {
					old[k].Value = ve.Value
				}
			}
		}
	}

	// 替换更新
	for k, v := range sets {
		for _, ve := range newSets {
			if v.Name == ve.Name {
				assign(sets[k].Content, ve.Content)
			}
		}
	}

	return sets, errmsg.SUCCES
}

// 从内存获取主题信息
func GetThemeInfo() (ThemeInfo, int) {
	infoByte := []byte(static.Get("theme_info"))

	var info ThemeInfo
	var code = errmsg.SUCCES

	if err := yaml.Unmarshal(infoByte, &info); err != nil {
		code = errmsg.ERROR
	}

	return info, code
}

// 保存设置项到内存
func saveThemeSetToStatic(data ThemeSetting) {
	for _, v := range data.Content {
		static.Set("theme_"+data.Name+"_"+v.Name, v.Value)
	}
}

// 保存设置项到数据库
func saveThemeSetToDB(data ThemeSetting) int {
	var set Settings
	set.Name = data.Name
	set.Role = "theme_" + getInfo().Name

	bytes, _ := json.Marshal(data.Content)
	set.Content = string(bytes)

	if err := db.Where("name = ? AND role = ?", set.Name, set.Role).Save(&set).Error; err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCES
}

// 修改设置项
func EditThemeSetting(data ThemeSetting) int {
	// 写入数据库
	if saveThemeSetToDB(data) != errmsg.SUCCES {
		return errmsg.ERROR
	}
	// 写入内存
	saveThemeSetToStatic(data)
	return errmsg.SUCCES
}
