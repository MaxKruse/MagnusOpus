package util

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetRequestFilter(c *fiber.Ctx) structs.RequestFilter {
	limit := c.Query("limit", "50")
	offset := c.Query("offset", "0")

	// convert limit and offset to ints
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 50
	}

	offsetInt, er := strconv.Atoi(offset)
	if er != nil {
		offsetInt = 0
	}

	return structs.RequestFilter{
		Limit:  limitInt,
		Offset: offsetInt,
	}
}

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

func CheckAuth(c *fiber.Ctx) (string, error) {
	sess, err := globals.SessionStore.Get(c)
	if err != nil {
		return "", err
	}

	// get authorization header
	auth := c.Get("Authorization")
	token := ""

	if len(auth) > 0 {
		token = auth[7:]
	} else {
		token = sess.ID()
	}

	// check if token is in database
	if !checkToken(token) {
		return "", fiber.ErrUnauthorized
	}

	if err := sess.Save(); err != nil {
		return "", err
	}

	return token, nil
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

func GetSelf(c *fiber.Ctx) (structs.User, error) {
	token, err := CheckAuth(c)
	if err != nil {
		return structs.User{}, err
	}

	user, err := GetUserFromSession(token)

	if err != nil {
		return structs.User{}, err
	}

	return user, nil
}
