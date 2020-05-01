package src

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"time"
)

type User struct {
	Id        int       `gorm:"type:int AUTO_INCREMENT unique"`
	Uid       string    `gorm:"type:varchar(26);primary_key"`
	Name      string    `gorm:"type:varchar(30);not null"`
	Email     string    `gorm:"type:varchar(254);not null"`
	Password  string    `gorm:"type:varchar(100);not null"`
	LastLogin time.Time `gorm:"not null"`
	CreateAt  time.Time `gorm:"not null"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Uid, validation.Required, validation.Length(26, 26)),
		validation.Field(&u.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required))
}
