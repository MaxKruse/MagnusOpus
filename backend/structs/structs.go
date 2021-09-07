package structs

import "gorm.io/gorm"

type Config struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_URL      string
}

// User struct
type Session struct {
	gorm.Model
	UserID       uint
	SessionToken string
	AccessToken  string `gorm:"type:varchar(64);"`
}

type User struct {
	gorm.Model
	RippleId int
	Session  Session
}

type RipplePing struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	UserId  int    `json:"user_id"`
}
