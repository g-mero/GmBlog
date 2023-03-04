package model

import (
	"fmt"
	"gmeroblog/utils/errmsg"

	"gorm.io/gorm"
)

type Category struct {
	ID       int    `gorm:"primary_key;auto_increment" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Slug     string `gorm:"not null;unique" json:"slug"`
	GroupId  int    `gorm:"not null" json:"group_id"`
	Desc     string `json:"desc"`
	CoverImg string `json:"cover_img"`
	Role     int    `gorm:"not null;" json:"role"`
}

type Catetree struct {
	Category
	Total    int64      `json:"total"`
	Children []Catetree `json:"children"`
}

// 初始化分类
func InitCate() int {
	var tmp Category
	err := db.Where("id = ?", 1).First(&tmp).Error
	if err == gorm.ErrRecordNotFound {
		tmp.ID = 1
		tmp.Name = "默认"
		tmp.Slug = "default"
		tmp.GroupId = 0
		tmp.Desc = "这是默认分类，它可以被修改但不能被删除"
		tmp.Role = 1
		return CreateCate(&tmp)
	} else {
		return errmsg.SUCCES
	}
}

// is category exist
func CheckCategory(category Category) (code int) {
	var id int
	err := db.Model(&Category{}).Select("id").Where("slug = ?", category.Slug).First(&id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errmsg.ERROR
	}
	if id > 0 && id != category.ID {
		return errmsg.ERROR_CATENAME_USED
	}
	return errmsg.SUCCES
}

// create category
func CreateCate(data *Category) int {
	fmt.Println(data)
	err := db.Create(&data).Error
	fmt.Println(data)
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

func GetSingleCate(id int) (Category, int) {
	var cate Category
	err = db.Where("id = ?", id).First(&cate).Error
	if err != nil {
		return cate, errmsg.ERROR_CATE_NOT_EXIST
	}
	return cate, errmsg.SUCCES
}

func GetCateBySlug(slug string) (Category, int) {
	var cate Category
	err = db.Where("slug = ?", slug).First(&cate).Error
	if err != nil {
		return cate, errmsg.ERROR_CATE_NOT_EXIST
	}
	return cate, errmsg.SUCCES
}

// get Category list
func GetCate() []Category {
	var cate []Category
	err = db.Where("role = ?", 0).Find(&cate).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	return cate
}

// get Category tree
func GetCates() []Catetree {
	var cates []Catetree
	var cateList []Category
	err = db.Order("id DESC").Find(&cateList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	for _, value := range cateList {
		if value.Role == 1 {
			var tmp Catetree
			tmp.Category = value
			tmp.Children = findChild(tmp, cateList)
			tmp.Total = 0
			for _, child := range tmp.Children {
				tmp.Total += child.Total
			}
			cates = append(cates, tmp)
		} else if value.GroupId == 0 {
			var tmp Catetree
			tmp.Category = value
			tmp.Total = getTotalOfCate(tmp.ID)
			cates = append(cates, tmp)
		}
	}
	return cates
}

// get cate's child
func findChild(cate Catetree, cateList []Category) []Catetree {
	var cates []Catetree
	for _, value := range cateList {
		if value.GroupId == cate.ID {
			var tmp Catetree
			tmp.Category = value
			tmp.Total = getTotalOfCate(tmp.ID)
			cates = append(cates, tmp)
		}
	}
	return cates
}

// 计算目录总文章数目
func getTotalOfCate(id int) int64 {
	var total int64
	err := db.Model(&Article{}).Where("cid = ?", id).Count(&total).Error
	if err != nil {
		return 0
	}
	return total
}

// edit Category
func EditCate(id int, data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"] = data.Name
	maps["slug"] = data.Slug
	maps["desc"] = data.Desc
	maps["group_id"] = data.GroupId
	maps["cover_img"] = data.CoverImg
	err = db.Model(&cate).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES

}

// delete Category
func DeleteCate(id int) int {
	var cate Category
	// 确保默认分类不会被删除
	if id == 1 {
		return errmsg.ERROR
	}
	err = db.Where("id = ?", id).Delete(&cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}
