package utils

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetRequestFilter(c *fiber.Ctx) structs.RequestFilter {
	defer TimeTrack(time.Now(), "GetRequestFilter")
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

func CheckAuth(c *fiber.Ctx) (string, error) {
	defer TimeTrack(time.Now(), "CheckAuth")
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

	return token, nil
}

func GetSelf(c *fiber.Ctx) (structs.User, error) {
	defer TimeTrack(time.Now(), "GetSelf")
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
