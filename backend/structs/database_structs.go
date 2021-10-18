package structs

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	JsonModel
	UserId       uint   `gorm:"not null" json:"-"`
	SessionToken string `gorm:"type:text" json:",omitempty"`
	AccessToken  string `gorm:"type:text" json:",omitempty"`
	RefreshToken string `gorm:"type:text" json:",omitempty"`
}

type User struct {
	JsonModel
	RippleId int       `json:"ripple_id,omitempty" gorm:"unique"`
	Username string    `json:"username,omitempty" gorm:"unique"`
	Sessions []Session `json:"sessions,omitempty"`
}

type Round struct {
	JsonModel
	TournamentId int       `json:"tournament_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	Active       bool      `json:"active"`
	DownloadPath string    `json:"download_path,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
}

type Staff struct {
	JsonModel
	TournamentId uint   `json:"tournament_id,omitempty"`
	User         *User  `json:"user,omitempty"`
	UserId       uint   `json:"-"`
	Role         string `json:"role"`
}

type Tournament struct {
	JsonModel
	Name                  string    `json:"name,omitempty" gorm:"unique"`
	Description           string    `json:"description,omitempty"`
	StartTime             time.Time `json:"start_time,omitempty"`
	EndTime               time.Time `json:"end_time,omitempty"`
	RegistrationStartTime time.Time `json:"registration_start_time,omitempty"`
	RegistrationEndTime   time.Time `json:"registration_end_time,omitempty"`
	Rounds                []Round   `json:"rounds,omitempty"`
	Staffs                []Staff   `json:"staffs,omitempty" gorm:"many2many:tournament_staff"`
	Registrations         []User    `json:"registrations,omitempty" gorm:"many2many:tournament_registrations"`
	Visible               bool      `json:"-"`
}

type BeatmapSubmittion struct {
	JsonModel
	Round        *Round `json:"round,omitempty"`
	RoundId      uint   `json:"-"`
	User         *User  `json:"user,omitempty"`
	UserId       uint   `json:"-"`
	Hash         string `json:"hash,omitempty"`
	DownloadPath string `json:"download_path,omitempty"`
	ToUse        bool   `json:"to_use"`
}

type JsonModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
