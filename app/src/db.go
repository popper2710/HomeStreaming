package src

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/my-repo/home_streaming/config"
	"math"
	"time"
)

type Secret struct {
}

func SqlConnect() *gorm.DB {
	var secret config.Secret
	secret.Init()
	DBMS := secret.Database.Image
	USER := secret.Database.User
	PASSWORD := secret.Database.Password
	HOST := "mysql"
	DBNAME := secret.Database.Dbname

	CONNECT := USER + ":" + PASSWORD + "@(" + HOST + ")/" + DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(DBMS, CONNECT)
		if err == nil {
			break
		} else {
			fmt.Println(err)
			fmt.Println("Continuing connect to DB..")

			// wait for container up
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
		}
	}
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Video{})
	return db
}
