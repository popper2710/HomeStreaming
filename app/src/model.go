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
	UpdatedAt time.Time `gorm:"not null"`
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
	Uid       string    `gorm:"type:varchar(36);not null;unique"`
	Name      string    `gorm:"type:varchar(500);not null;default:'Untitled'"`
	IsEncode  bool      `gorm:"type:bool;not null;default:false"`
	UpdatedAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UserId    int       `gorm:"type:int"`
	User      User
}

func (v Video) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Uid, validation.Required, is.UUIDv4),
	)
}
