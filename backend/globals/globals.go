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
		Session: &structs.Session{
			SessionToken: bearer,
		},
	}

	DBConn.Preload("Session").First(&user, user)
	Logger.WithFields(logrus.Fields{
		"token": bearer,
	}).Debug("Checking token")

	return user.ID != 0
}

func CheckAuth(c *fiber.Ctx) (string, error) {
	// get authorization header
	auth := c.Get("Authorization")
	Logger.WithFields(logrus.Fields{
		"auth": auth,
	}).Debug("Checking token")
	if auth == "" {
		return "", fiber.ErrUnauthorized
	}

	// get token
	token := auth[7:]
	Logger.WithFields(logrus.Fields{
		"token": token,
	}).Debug("Checking token")

	// check if token is in database
	if !checkToken(token) {
		return "", fiber.ErrUnauthorized
	}

	return token, nil
}

func GetSelf(c *fiber.Ctx) (structs.User, error) {
	token, err := CheckAuth(c)
	if err != nil {
		return structs.User{}, err
	}

	search := structs.User{
		Session: &structs.Session{
			SessionToken: token,
		},
	}

	res := structs.User{}
	err = DBConn.Find(&res, search).Error

	if err != nil {
		return structs.User{}, err
	}

	return res, nil
}
