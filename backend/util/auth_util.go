package util

import (
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func checkToken(bearer string) bool {
	// check if token is in database
	user := structs.User{}
	sess := structs.Session{SessionToken: bearer}
	localDB := globals.DBConn

	globals.Logger.WithField("sess", sess).Debug("Checking token from database")
	localDB.Find(&sess, sess)
	globals.Logger.WithField("sess", sess).Debug("Checking token from database")

	// get user from sess
	localDB.Preload("Session").First(&user, "session_id = ?", sess.ID)
	globals.Logger.WithField("user", user).Debug("User found")
	return user.ID != 0
}

func GetUserFromSession(sessionToken string) (structs.User, error) {
	user := structs.User{}
	sess := structs.Session{SessionToken: sessionToken}
	localDB := globals.DBConn

	err := localDB.First(&sess, sess).Error

	if err != nil {
		return structs.User{}, err
	}

	err = localDB.Preload("Session").First(&user, "session_id = ?", sess.ID).Error
	if err != nil {
		return structs.User{}, err
	}

	return user, nil
}
