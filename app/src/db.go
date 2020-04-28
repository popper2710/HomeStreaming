package src

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"time"
)

type Secret struct {
	Database struct {
		Image    string `yaml:"image"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}
}

func SqlConnect() *gorm.DB {
	buf, err := ioutil.ReadFile("./config/secret.yaml")
	if err != nil {
		panic(err)
	}
	var secret Secret
	err = yaml.Unmarshal(buf, &secret)
	if err != nil {
		panic(err)
	}

	DBMS := secret.Database.Image
	USER := secret.Database.User
	PASSWORD := secret.Database.Password
	HOST := "mysql"
	DBNAME := secret.Database.Dbname

	CONNECT := USER + ":" + PASSWORD + "@(" + HOST + ")/" + DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	var db *gorm.DB
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
	return db
}
