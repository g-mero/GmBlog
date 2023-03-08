package model

import (
	"gmeroblog/utils/errmsg"
	"log"
	"math/rand"
	"time"
)

type Mood struct {
	ID        int       `gorm:"primarykey" json:"id"`
	Content   string    `gorm:"not null" json:"content"`
	Private   bool      `gorm:"DEFAULT:false" json:"private"`
	CreatedAt time.Time `json:"created_at"`
}

// 添加一条心情
func AddMood(data *Mood) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// 获取心情
func GetMoods(pageSize int, pageNum int) ([]Mood, int, int64) {
	var moods []Mood
	var total int64
	if err := db.Model(&moods).Order("id DESC").Count(&total).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&moods).Error; err != nil {
		return nil, errmsg.ERROR, total
	}
	return moods, errmsg.SUCCES, total
}

var errMoods = []string{"羡慕那些坚定的人", "好好学习，天天运动", "人应该有梦想"}

// 获取一条随机心情的内容，只会获取对外开放的心情
func GetMoodRand() string {
	var mood Mood
	var max int64
	db.Model(&mood).Count(&max)
	if max <= 0 {
		return errMoods[rand.Intn(3)]
	}

	var randID = rand.Intn(int(max)) + 1 //随机的id
	if err := db.Where("id = ? AND private = ?", randID, 0).First(&mood).Error; err != nil {
		log.Println("[MOOD]", err)
		return errMoods[rand.Intn(3)]
	}

	return mood.Content
}

// 删除心情
func DelMood(id int) int {
	err := db.Where("id = ?", id).Delete(&cate{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}
