package globals

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Logger *logrus.Logger
	Config structs.Config
	DBConn *gorm.DB
)

func checkToken(bearer string) bool {

	// check if token is in database
	user := structs.User{
		Session: structs.Session{
			SessionToken: bearer,
		},
	}

	DBConn.Preload("Session").First(&user, user)
	Logger.WithFields(logrus.Fields{
		"token": bearer,
	}).Debug("Checking token")

	return user.ID != 0
}

func CheckAuth(c *fiber.Ctx) error {
	// get authorization header
	auth := c.Get("Authorization")
	Logger.WithFields(logrus.Fields{
		"auth": auth,
	}).Debug("Checking token")
	if auth == "" {
		return fiber.ErrUnauthorized
	}

	// get token
	token := auth[7:]
	Logger.WithFields(logrus.Fields{
		"token": token,
	}).Debug("Checking token")

	// check if token is in database
	if !checkToken(token) {
		return fiber.ErrUnauthorized
	}

	return nil
}
