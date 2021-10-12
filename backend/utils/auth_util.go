package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
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

func CheckAuth(token string) (string, error) {

	// check if token is in database
	if !checkToken(token) {
		return "", fiber.ErrUnauthorized
	}

	return token, nil
}

func GetUserFromSession(sessionToken string) (structs.User, error) {
	// check if token is in database
	user := structs.User{}
	sess := structs.Session{SessionToken: sessionToken}
	localDB := globals.DBConn

	err := localDB.Find(&sess, sess).Error
	if err != nil {
		return user, errors.New("invalid session token")
	}

	// get session from db
	err = localDB.First(&sess, sess).Error
	if err != nil {
		return user, errors.New("invalid session token " + sessionToken)
	}

	// get user from sess
	localDB.Preload("Sessions").First(&user, sess.UserId)

	return user, nil
}
