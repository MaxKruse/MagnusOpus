package util

import (
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func checkToken(bearer string) bool {
	// check if token is in database
	user := structs.User{}
	sess := structs.Session{SessionToken: bearer}

	globals.Logger.WithField("sess", sess).Debug("Checking token from database")
	globals.DBConn.Debug().Find(&sess, sess)
	globals.Logger.WithField("sess", sess).Debug("Checking token from database")

	// get user from sess
	globals.DBConn.Debug().Preload("Session").First(&user, "session_id = ?", sess.ID)
	globals.Logger.WithField("user", user).Debug("User found")
	return user.ID != 0
}

func GetUserFromSession(sessionToken string) (structs.User, error) {
	user := structs.User{}
	sess := structs.Session{SessionToken: sessionToken}

	err := globals.DBConn.First(&sess, sess).Debug().Error

	if err != nil {
		return structs.User{}, err
	}

	err = globals.DBConn.Debug().Preload("Session").First(&user, "session_id = ?", sess.ID).Error
	if err != nil {
		return structs.User{}, err
	}

	return user, nil
}
