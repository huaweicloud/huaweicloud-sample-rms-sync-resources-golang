package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"huaweicloud-sample-rms-sync-resources/models"
	"io/ioutil"
	"log"
)

var DB *gorm.DB
var SQL_TEMPLATE_UPSERT_ONE_RESOURCE string

func init() {
	logrus.Info("Database initialization..")
	DB = GetDatabase()
	err1 := migrateTables()
	if err1 != nil {
		log.Fatal(err1)
	}
	SQL_TEMPLATE_UPSERT_ONE_RESOURCE = readSql()
}

func migrateTables() error {
	if err := migrateResourceTable(); err != nil {
		return err
	}
	if err := migrateTagTable(); err != nil {
		return err
	}
	return nil
}

func migrateResourceTable() error {
	err := DB.AutoMigrate(&models.Resource{})
	if err != nil {
		return err
	}
	return nil
}

func migrateTagTable() error {
	err := DB.AutoMigrate(&models.Tag{})
	if err != nil {
		return err
	}
	return nil
}

func readSql() string {
	f, err := ioutil.ReadFile("conf/upsert_one_resource.sql")
	if err != nil {
		log.Fatalf("Failed to read sql file, err: %v", err)
	}
	return string(f)
}
