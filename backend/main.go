package main

import (
	// Import logrus

	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/routes"
	"github.com/maxkruse/magnusopus/backend/routes/tournaments"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

func init() {
	// Setup global state before main

	// Create a new instance of the logger
	globals.Logger = logrus.New()

	// Set the log level
	globals.Logger.SetLevel(logrus.DebugLevel)

	// Set log config
	globals.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.1234",
		PrettyPrint:     true,
		DataKey:         "data",
	})

	// Get POSTGRES_ env variables
	globals.Config.POSTGRES_USER = os.Getenv("POSTGRES_USER")
	globals.Config.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	globals.Config.POSTGRES_DB = os.Getenv("POSTGRES_DB")
	globals.Config.POSTGRES_URL = os.Getenv("POSTGRES_URL")

	// connect to database
	err := error(nil)
	globals.DBConn, err = gorm.Open(postgres.Open(globals.Config.POSTGRES_URL), &gorm.Config{
		PrepareStmt: true,
		Logger:      glog.Default.LogMode(glog.Info),
	})
	if err != nil {
		globals.Logger.Fatal(err)
		os.Exit(1)
	}
	globals.Logger.Debug("Connected to database")

	// Migrate Tables
	globals.DBConn.AutoMigrate(&structs.User{})
	globals.DBConn.AutoMigrate(&structs.Session{})
	globals.DBConn.AutoMigrate(&structs.Round{})
	globals.DBConn.AutoMigrate(&structs.Tournament{})
	globals.Logger.Debug("Migrated")

	globals.SessionStore = session.New(session.Config{
		Expiration: time.Hour * 24 * 7 * 31, // 1 Month
		Storage: sqlite3.New(sqlite3.Config{
			Reset:    false,
			Database: "/storage/sessions.db",
		}),
	})

	globals.Logger.Info("Starting Magnusopus backend")
	globals.Logger.WithFields(logrus.Fields{"config": globals.Config}).Debug("Config")
}

func checkSessionCookie(c *fiber.Ctx) error {
	// check if session is valid
	_, err := globals.SessionStore.Get(c)
	if err != nil {
		globals.Logger.WithFields(logrus.Fields{"error": err}).Error("Session error")
		return err
	}

	// Check if session cookie is set
	session_token, err := globals.CheckAuth(c)

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// check if auth_token is in database
	user := structs.NewUser()
	user.Session.SessionToken = session_token

	globals.DBConn.Preload("Session").First(&user, user)
	if user.Session.ID == 0 {

		// no valid session found with given credentials, return error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "no session found, please login",
		})
	}

	return c.Next()
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork: false, // true = multithreaded, false = singlethreaded
	})

	// use middlewares
	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05.1234",
	}))
	app.Use(etag.New())
	app.Use(compress.New())

	// oauth routes
	oauth := app.Group("/oauth")
	oauth.Get("/ripple", routes.GetOAuthRipple)

	api := app.Group("/api")

	// use custom middleware
	api.Use(checkSessionCookie)

	v1 := api.Group("/v1")

	v1.Get("/me", routes.Me)

	v1.Get("/tournaments", tournaments.GetTournaments)
	v1.Get("/tournaments/:id", tournaments.GetTournament)

	app.Listen(":5000")
}
