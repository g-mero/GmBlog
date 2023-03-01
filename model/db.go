package model

import (
	"fmt"
	"gmeroblog/utils/config"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var err error

func con_sqlite3() (*gorm.DB, error) {
	sdb, err := gorm.Open(sqlite.Open(config.DbFile), &gorm.Config{
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	return sdb, err
}

func InitDb() {
	db, err = con_sqlite3()
	if err != nil {
		fmt.Println("can not connect to sqlite3", err)
		os.Exit(1)
	}

	// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
	// _ = db.AutoMigrate(&User{}, &Article{}, &Category{}, &Diary{}, &Settings{}, &Comment{})

	// 数据库基本内容的初始化(理论上只需要运行一次)
	if InitCate()+InitSet()+InitUser() != 600 {
		fmt.Println("基础内容初始化失败")
		os.Exit(1)
	}

	sqlDB, _ := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}
