package model

import (
	"gmeroblog/utils/errmsg"
	"time"
)

type Diary struct {
	Time      string    `gorm:"not null;primaryKey" json:"time"`
	Desc      string    `json:"desc"`
	Content   string    `gorm:"not null" json:"content"`
	Extra     string    `json:"extra"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// create diary
func CreateDiary(data *Diary) int {
	if data.Time == "" {
		now := time.Now()
		data.Time = now.Format("2006-01-02")
	}
	if data.Desc == "" {
		rs := []rune(data.Content)
		if len(rs) > 40 {
			data.Desc = string(rs[0:40])
		} else {
			data.Desc = data.Content
		}
	}
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// get single diary
func GetSingleDiary(time string) (Diary, int) {
	var art Diary
	if time == "lastest" {
		err = db.Last(&art).Error
	} else {
		err = db.Where("time = ?", time).First(&art).Error
	}
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCES
}

// get diary list
func GetDiary(opt int, value string) ([]Diary, int, int64) {
	var diary []Diary
	var total int64
	if opt == 0 {
		err = db.Model(&diary).Count(&total).Find(&diary).Error
	} else if opt == 1 {
		err = db.Model(&diary).Where("time like ?", value+"-__").Count(&total).Find(&diary).Error
	} else {
		err = db.Model(&diary).Where("time like ?", value+"-__-__").Count(&total).Find(&diary).Error
	}

	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return diary, errmsg.SUCCES, total
}

// get diary date by years or month
func GetDiaryDate(opt int, value string) ([]string, int, int64) {
	var diaryDate []string
	var total int64
	if opt == 0 {
		err = db.Model(&Diary{}).Select("time").Count(&total).Find(&diaryDate).Error
	} else if opt == 1 {
		err = db.Model(&Diary{}).Select("time").Where("time like ?", value+"-__").Count(&total).Find(&diaryDate).Error
	} else {
		err = db.Model(&Diary{}).Select("time").Where("time like ?", value+"-__-__").Count(&total).Find(&diaryDate).Error
	}

	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return diaryDate, errmsg.SUCCES, total
}

// edit Diary
func EditDiary(time string, data *Diary) int {
	var art Diary
	var maps = make(map[string]interface{})
	maps["time"] = data.Time
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["extra"] = data.Extra
	err = db.Model(&art).Where("time = ?", time).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES

}

// delete Diary
func DeleteDiary(time string) int {
	var art Diary
	err = db.Where("time = ?", time).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}
