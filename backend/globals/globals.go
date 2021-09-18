package globals

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Logger       *logrus.Logger
	Config       structs.Config
	DBConn       *gorm.DB
	SessionStore *session.Store
)

func checkToken(bearer string) bool {
	// check if token is in database
	user := structs.NewUser()
	user.Session.SessionToken = bearer

	DBConn.Debug().Preload("Session").First(&user, user)
	return user.ID != 0
}

func CheckAuth(c *fiber.Ctx) (string, error) {
	sess, err := SessionStore.Get(c)
	if err != nil {
		return "", err
	}

	// get authorization header
	auth := c.Get("Authorization")
	token := ""

	if len(auth) > 0 {
		token := auth[7:]
		Logger.WithField("token", token).Debug("Checking token from Authorization")
	} else {
		token = sess.ID()
		Logger.WithField("token", token).Debug("Checking token from session_id")
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

func GetSelf(c *fiber.Ctx) (structs.User, error) {
	token, err := CheckAuth(c)
	if err != nil {
		return structs.NewUser(), err
	}

	search := structs.NewUser()
	search.Session.SessionToken = token

	Logger.WithField("token", token).Debug("Getting user from search")

	res := structs.NewUser()
	err = DBConn.Debug().Preload("Session").Find(&res, search).Error

	if err != nil {
		return structs.NewUser(), err
	}

	return res, nil
}
