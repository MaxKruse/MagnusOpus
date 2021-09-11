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
	gorm.Model
	UserID       uint
	SessionToken string `gorm:"type:text"`
	AccessToken  string `gorm:"type:text"`
	RefreshToken string `gorm:"type:text"`
}

type User struct {
	gorm.Model
	RippleId int      `json:"ripple_id ,omitempty"`
	BanchoId int      `json:"bancho_id ,omitempty"`
	Username string   `json:"username ,omitempty"`
	Session  *Session `json:"session ,omitempty"`
}

type RippleSelf struct {
	UserId   int    `json:"id ,omitempty"`
	Username string `json:"username ,omitempty"`
}

type BanchoMe struct {
	Id       int    `json:"id ,omitempty"`
	Username string `json:"username ,omitempty"`
}

type Round struct {
	gorm.Model
	TournamentId int       `json:"tournament_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
}

type Staff struct {
	gorm.Model
	UserId int
	Role   string `json:"role"`
}

type Tournament struct {
	gorm.Model
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	DownloadPath string    `json:"file,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
	rounds       *[]Round  `json:"rounds,omitempty"`
	staffs       *[]Staff  `json:"staffs,omitempty"`
}
