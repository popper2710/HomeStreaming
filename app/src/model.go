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
	CreatedAt time.Time `gorm:"not null"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Uid, validation.Required, validation.Length(26, 26)),
		validation.Field(&u.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required))
}

type Video struct {
	Id        int       `gorm:"type:int AUTO_INCREMENT"`
	Uid       string    `gorm:"type:string;varchar(36);not null;unique"`
	Name      string    `gorm:"type:varchar(500)';not null;default:'Untitled'"`
	CreatedAt time.Time `gorm:"not null"`
}

func (v Video) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Uid, validation.Required, is.UUIDv4),
	)
}
