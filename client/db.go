package client

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"picasso/config"
	"time"
)

var Orm *gorm.DB

func InitDB() {
	Orm = gormDB()
}

func gormDB() *gorm.DB {
	dsn := config.SysYamlconfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	//err = db.AutoMigrate(&User{}, &Question{}, &Answer{})
	//if err != nil {
	//	log.Fatal(err)
	//}

	mysqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetMaxOpenConns(10)
	mysqlDB.SetConnMaxLifetime(time.Second * 30)
	return db
}
