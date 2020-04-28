package src

import "time"

type User struct {
	Id        int       `gorm:"type:int(10);primary_key;auto_increment"`
	Name      string    `gorm:"type:varchar(30);not null"`
	Password  string    `gorm:"type:varchar(100);not null"`
	LastLogin time.Time `gorm:"not null"`
	CreateAt  time.Time `gorm:"not null"`
}
