package utils

import (
	"errors"
	"log"

	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetTournament(tournament_id uint) (structs.Tournament, error) {
	localDB := globals.DBConn

	t := structs.Tournament{}
	err := localDB.Preload("Staffs.User").First(&t, tournament_id).Error
	if err != nil {
		return t, err
	}

	return t, nil
}

func CanViewTournament(user_id uint, tournament_id uint) error {
	t, err := GetTournament(tournament_id)
	if err != nil {
		return err
	}

	for _, staff := range t.Staffs {
		if staff.User.ID == user_id {
			switch staff.Role {
			case "owner", "admin", "mod", "judge":
				return nil
			}
		}
	}

	if t.Visible {
		return nil
	}

	return errors.New("not allowed to view tournament: need owner, admin, mod, judge")
}

func CanEditTournament(user_id uint, tournament_id uint) error {
	localDB := globals.DBConn

	t := structs.Tournament{}
	err := localDB.Preload("Staffs.User").First(&t, tournament_id).Error
	if err != nil {
		return err
	}

	for _, staff := range t.Staffs {
		if staff.User.ID == user_id {
			switch staff.Role {
			case "owner", "admin":
				return nil
			}
		}
	}

	return errors.New("not allowed to edit tournament: need owner, admin")
}

func CanEditRounds(user_id uint, tournament_id uint) error {
	localDB := globals.DBConn

	t := structs.Tournament{}
	err := localDB.Preload("Staffs.User").First(&t, tournament_id).Error
	if err != nil {
		return err
	}

	for _, staff := range t.Staffs {
		if staff.User.ID == user_id {
			switch staff.Role {
			case "owner", "admin", "mod":
				return nil
			}
		}
	}

	return errors.New("not allowed to edit tournament: need owner, admin, mod")
}

func CanJudge(user_id uint, tournament_id uint) error {
	localDB := globals.DBConn

	t := structs.Tournament{}
	err := localDB.Preload("Staffs.User").First(&t, tournament_id).Error
	if err != nil {
		return err
	}

	for _, staff := range t.Staffs {
		if staff.User.ID == user_id {
			log.Println("User Match:", staff.Role)
			switch staff.Role {
			case "owner", "judge":
				return nil
			}
		}
	}

	return errors.New("not allowed to edit tournament: need owner, judge")
}
