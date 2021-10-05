package utils

import (
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func checkToken(bearer string) bool {
	user, err := GetUserFromSession(bearer)
	if err != nil {
		return false
	}
	return user.ID != 0
}

func GetUserFromSession(sessionToken string) (structs.User, error) {
	// check if token is in database
	user := structs.User{}
	sess := structs.Session{SessionToken: sessionToken}
	localDB := globals.DBConn

	err := localDB.Find(&sess, sess).Error
	if err != nil {
		return user, err
	}
	globals.Logger.WithField("sess", sess).Debug("Got session from token")

	// get session from db
	err = localDB.First(&sess, sess).Error
	if err != nil {
		return user, err
	}

	// get user from sess
	localDB.Preload("Sessions").First(&user, sess.UserId)
	globals.Logger.WithField("user", user).Debug("User found")

	return user, nil
}
