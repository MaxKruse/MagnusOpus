package structs

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	JsonModel
	SessionToken string `gorm:"type:text" json:",omitempty"`
	AccessToken  string `gorm:"type:text" json:",omitempty"`
	RefreshToken string `gorm:"type:text" json:",omitempty"`
}

type User struct {
	JsonModel
	RippleId  int     `json:"ripple_id ,omitempty" gorm:"unique"`
	Username  string  `json:"username ,omitempty" gorm:"unique"`
	Session   Session `json:"session ,omitempty"`
	SessionId uint    `json:"-"`
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
	TournamentId uint   `json:"tournament_id,omitempty"`
	UserId       uint   `json:"user_id"`
	Role         string `json:"role"`
}

type Tournament struct {
	JsonModel
	Name         string    `json:"name,omitempty" gorm:"unique"`
	Description  string    `json:"description,omitempty"`
	DownloadPath string    `json:"download_path,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
	Rounds       []Round   `json:"rounds,omitempty"`
	Staffs       []Staff   `json:"staffs,omitempty"`
	Visible      bool      `json:"-"`
}

type JsonModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
