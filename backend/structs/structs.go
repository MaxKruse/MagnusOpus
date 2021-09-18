package structs

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_URL      string
}

// User struct
type Session struct {
	JsonModel
	UserID       uint   `json:"-"`
	SessionToken string `gorm:"type:text" json:",omitempty"`
	AccessToken  string `gorm:"type:text" json:",omitempty"`
	RefreshToken string `gorm:"type:text" json:",omitempty"`
}

type User struct {
	JsonModel
	RippleId int      `json:"ripple_id ,omitempty"`
	Username string   `json:"username ,omitempty"`
	Session  *Session `json:"session ,omitempty"`
}

func NewUser() User {
	r := User{}
	r.Session = new(Session)

	return r
}

type RippleSelf struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	UserId   int    `json:"id"`
	Username string `json:"username"`
}

type BanchoMe struct {
	Id       int    `json:"id ,omitempty"`
	Username string `json:"username ,omitempty"`
}

type Round struct {
	JsonModel
	TournamentId int       `json:"tournament_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
}

type Staff struct {
	JsonModel
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}

type Tournament struct {
	JsonModel
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	DownloadPath string    `json:"file,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
	rounds       []*Round  `json:"rounds,omitempty"`
	staffs       []*Staff  `json:"staffs,omitempty"`
}

type JsonModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
