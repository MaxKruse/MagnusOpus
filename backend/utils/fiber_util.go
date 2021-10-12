package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/performance"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetRequestFilter(c *fiber.Ctx) structs.RequestFilter {
	defer performance.TimeTrack(time.Now(), "GetRequestFilter")
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

func GetSelfFromDB(c *fiber.Ctx) (structs.User, error) {
	token, err := GetTokenFromRequest(c)
	if err != nil {
		return structs.User{}, err
	}

	user, err := GetUserFromSession(token)

	if err != nil {
		return structs.User{}, err
	}

	return user, nil
}

func GetTokenFromRequest(c *fiber.Ctx) (string, error) {
	token := c.Get("Authorization")
	if token != "" {
		return token, nil
	}
	token = c.Cookies("session_id", "")
	if token != "" {
		return token, nil
	}

	return "", errors.New("no token in request")
}

func IsSuperadmin(c *fiber.Ctx) error {
	self, err := GetSelfFromDB(c)

	if err != nil {
		return err
	}

	// check if user.RippleId is in globals.Superadmins
	for _, superadmin := range globals.AllowedSuperadmin {
		if superadmin == self.RippleId {
			return nil
		}
	}

	return errors.New("user is not a superadmin")
}

func DefaultErrorMessage(c *fiber.Ctx, err error, code int) error {
	if err != nil {
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}

func StringToUint32(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}
