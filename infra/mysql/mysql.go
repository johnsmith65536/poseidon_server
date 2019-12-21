package mysql

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var db *gorm.DB

const Addr = "localhost:3306"

func Init() {
	var err error
	dsn := "root:mysql@tcp(" + Addr + ")/poseidon?charset=utf8mb4&parseTime=True&timeout=5000ms&readTimeout=5000ms&writeTimeout=5000ms&loc=Asia%2FChongqing&interpolateParams=true&maxAllowedPacket=0"
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.DB().SetMaxOpenConns(50)
	db.DB().SetMaxIdleConns(50)
	db.DB().SetConnMaxLifetime(200 * time.Second)
	db.SingularTable(true)
}

func Test() error {
	return nil
}
