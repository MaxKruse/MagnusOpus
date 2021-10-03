package structs

import (
	"errors"
	"strings"
	"time"

	"github.com/maxkruse/magnusopus/backend/performance"
	"gorm.io/gorm"
)

func (t *Tournament) ActivateRound(name string) error {
	defer performance.TimeTrack(time.Now(), "ActivateRound")

	if t == nil {
		return errors.New("activate round: tournament is nil")
	}

	if t.Rounds == nil {
		return errors.New("activate round: tournament.rounds is nil")
	}

	if len(t.Rounds) == 0 {
		return errors.New("activate round: tournament.rounds is empty")
	}

	changed := false

	for _, r := range t.Rounds {
		// Deactivate all rounds
		r.Active = false

		// Only activate if the name matches
		if strings.EqualFold(r.Name, name) {
			r.Active = true
			changed = true
		}
	}

	if !changed {
		return errors.New("activate round: no round with name " + name)
	}

	return nil
}

func (t Tournament) ValidTournament(localDB *gorm.DB) error {
	if t.Name == "" {
		return errors.New("name is required")
	}

	if t.Description == "" {
		return errors.New("description is required")
	}

	zeroTime := time.Time{}

	if t.StartTime == zeroTime {
		return errors.New("start_time is required (ISO 8601) (RFC 3339)")
	}

	if t.EndTime == zeroTime {
		return errors.New("end_time is required (ISO 8601) (RFC 3339)")
	}

	// Check if the time is in the future
	if t.StartTime.Before(time.Now()) {
		return errors.New("start_time must be in the future")
	}

	if t.EndTime.Before(time.Now()) {
		return errors.New("end_time must be in the future")
	}

	// check if end_time is at least 3 days after start_time
	if t.EndTime.Sub(t.StartTime) < (3 * 24 * time.Hour) {
		return errors.New("end_time must be at least 3 days after start_time")
	}

	res := Tournament{}
	localDB.Where(t).First(&res)

	if res.ID != 0 {
		return errors.New("name must be unique")
	}

	return nil
}

func (t Round) ValidRound(localDB *gorm.DB) error {
	if t.Name == "" {
		return errors.New("name is required")
	}

	if t.Description == "" {
		return errors.New("description is required")
	}

	if t.DownloadPath == "" {
		return errors.New("download_path is required")
	}

	zeroTime := time.Time{}

	if t.StartTime == zeroTime {
		return errors.New("start_time is required (ISO 8601) (RFC 3339)")
	}

	if t.EndTime == zeroTime {
		return errors.New("end_time is required (ISO 8601) (RFC 3339)")
	}

	// Check if the time is in the future
	if t.StartTime.Before(time.Now()) {
		return errors.New("start_time must be in the future")
	}

	if t.EndTime.Before(time.Now()) {
		return errors.New("end_time must be in the future")
	}

	// check if end_time is at least 3 days after start_time
	if t.EndTime.Sub(t.StartTime) < (3 * 24 * time.Hour) {
		return errors.New("end_time must be at least 3 days after start_time")
	}

	res := Round{}
	localDB.Where(t).First(&res)

	if res.ID != 0 {
		return errors.New("name must be unique")
	}

	return nil
}
