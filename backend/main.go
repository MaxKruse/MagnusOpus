package main

import (
	// Import logrus

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/routes"
	"github.com/maxkruse/magnusopus/backend/routes/tournaments"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// Setup global state before main

	// Create a new instance of the logger
	globals.Logger = logrus.New()

	// Set the log level
	globals.Logger.SetLevel(logrus.DebugLevel)

	// Set log config
	globals.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
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
	globals.DBConn, err = gorm.Open(postgres.Open(globals.Config.POSTGRES_URL), &gorm.Config{})
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

	globals.Logger.Info("Starting Magnusopus backend")
	globals.Logger.WithFields(logrus.Fields{"config": globals.Config}).Debug("Config")
}

func checkSessionCookie(c *fiber.Ctx) error {
	// Check if session cookie is set
	session_token := c.Cookies("session_token")
	if session_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "no session_token, please login",
		})
	}

	// check if auth_token is in database
	user := structs.User{}
	user.Session = &structs.Session{SessionToken: session_token}

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
	app.Use(logger.New())
	app.Use(etag.New())
	app.Use(compress.New())
	app.Use(recover.New())

	// oauth routes
	oauth := app.Group("/oauth")
	oauth.Get("/ripple", routes.GetOAuthRipple)
	oauth.Get("/bancho", routes.GetOAuthBancho)

	api := app.Group("/api")

	// use custom middleware
	api.Use(checkSessionCookie)

	v1 := api.Group("/v1")

	v1.Get("/me", routes.Me)

	v1.Get("/tournaments", tournaments.GetTournaments)
	v1.Get("/tournaments/:id", tournaments.GetTournament)

	app.Listen(":5000")
}
