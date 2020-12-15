package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"huaweicloud-sample-rms-sync-resources/config"
	"log"
	"sync"
	"time"
)

var singletonDB *gorm.DB
var once sync.Once

func GetDatabase() *gorm.DB {
	once.Do(InitDatabase)
	return singletonDB
}

func InitDatabase() {
	cfg := config.GetConfig()
	mysql_cfg := cfg.Mysql
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True",
		mysql_cfg.Username, mysql_cfg.Password, mysql_cfg.Network, mysql_cfg.Server, mysql_cfg.Port, mysql_cfg.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err1 := db.DB()
	if err1 != nil {
		log.Fatal(err1)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(100 * time.Second)
	singletonDB = db
}
