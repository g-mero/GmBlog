package model

import (
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

// api输出类型
type OutThemeSetting struct {
	ID int `json:"id"`
	ThemeSetting
}

type oneSet struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
	Label string `yaml:"label" json:"label"`
	Desc  string `yaml:"desc" json:"desc"`
	Type  string `yaml:"type" json:"type"` // text | textarea | radio | select | switch
}

var ThemeName = getInfo().Name
var themeRole = "Theme_" + ThemeName
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

	sets := getThemeSetFormFile()

	if sets == nil {
		log.Fatal("[Theme]", "获取设置项失败")
	}

	if initThemeDb(sets) != errmsg.SUCCES {
		log.Fatal("[Theme]", "初始化数据库失败")
	}

	updateThemeSetWithDB(sets)
}

// 初始化数据库
func initThemeDb(sets []ThemeSetting) int {
	var setInit = func(set ThemeSetting) int {
		var tmp Settings
		err := db.Where("name = ? AND role = ?", set.Name, themeRole).First(&tmp).Error
		if err == gorm.ErrRecordNotFound {
			dbSet := themeSetToSettings(set)
			err = db.Create(&dbSet).Error
			if err != nil {
				return errmsg.ERROR
			}
		}
		return errmsg.SUCCES
	}

	for _, ts := range sets {
		if setInit(ts) == errmsg.ERROR {
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

// 从文件中获取主题设置项
func getThemeSetFormFile() []ThemeSetting {
	// 文件中读取到的设置项
	var sets []ThemeSetting

	file, err := os.ReadFile(setPath)
	if err != nil {
		return nil
	}

	// 从文件中读取数据(此时value还是默认值，需要与数据库进行比对)
	if err := yaml.Unmarshal(file, &sets); err != nil {
		return nil
	}

	return sets
}

// 将ThemeSetting类型转化为Settings类型
func themeSetToSettings(themeSet ThemeSetting) Settings {
	var set Settings

	set.ID = 0
	set.Name = themeSet.Name
	set.Role = themeRole
	set.Content = jsonString(oneSetsToMap(themeSet.Content))
	return set

}

// 将oneSet数组转Map
func oneSetsToMap(sets []oneSet) map[string]string {
	setMap := make(map[string]string, len(sets))
	for _, v := range sets {
		setMap[v.Name] = v.Value
	}
	return setMap
}

// 将文件设置项与数据库中进行比对更新, 并写入内存
func updateThemeSetWithDB(fileSets []ThemeSetting) ([]OutThemeSetting, int) {
	// 从数据库中读取设置项
	dbSets, code := GetSettings(themeRole)
	if code != errmsg.SUCCES {
		// 数据库读取失败
		return nil, errmsg.ERROR
	}

	// 将数据进行比对更新，注意这里会直接修改OLD
	var assign = func(old []oneSet, new map[string]string) {
		for k, v := range old {
			for ke, ve := range new {
				if v.Name == ke {
					old[k].Value = ve
				}
			}
		}
	}

	var outThemeSettings []OutThemeSetting

	// 开始比对
	for k, v := range fileSets {
		for _, ve := range dbSets {
			if v.Name == ve.Name {
				assign(fileSets[k].Content, ve.Content)
				var tmp OutThemeSetting
				tmp.ThemeSetting = fileSets[k]
				saveThemeSetToStatic(ve)
				tmp.ID = ve.ID
				// 转化成out格式
				outThemeSettings = append(outThemeSettings, tmp)
			}
		}
	}

	return outThemeSettings, errmsg.SUCCES
}

// 获取主题的所有设置项
func GetThemeSettings() ([]OutThemeSetting, int) {
	// 文件中读取到的设置项
	sets := getThemeSetFormFile()

	if sets == nil {
		return nil, errmsg.ERROR
	}

	return updateThemeSetWithDB(sets)
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
func saveThemeSetToStatic(data OutSettings) {
	for k, v := range data.Content {
		static.Set("THEME_"+ThemeName+"_"+data.Name+"_"+k, v)
	}
}

// 修改设置项
func EditThemeSetting(data OutSettings) int {
	// 写入数据库
	var set Settings
	if err := db.Where("id = ?", data.ID).First(&set).Error; err != nil || set.Name != data.Name || set.Role != themeRole {
		return errmsg.ERROR
	}
	set.Content = jsonString(data.Content)
	err := db.Save(&set).Error
	if err != nil {
		return errmsg.ERROR
	}
	saveThemeSetToStatic(data)
	return errmsg.SUCCES
}
