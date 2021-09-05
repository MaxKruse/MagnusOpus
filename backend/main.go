package main

import (
	// Import logrus

	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/routes"
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
	globals.Logger.Debug("Migrated")

	globals.Logger.Info("Starting Magnusopus backend")
	globals.Logger.WithFields(logrus.Fields{"config": globals.Config}).Debug("Config")
}

func checkSessionCookie(c *fiber.Ctx) error {
	// Check if session cookie is set
	accessToken := c.Cookies("ripple_token")
	if accessToken == "" {
		// Redirect to login
		return c.Redirect("/oauth")
	}

	// check if auth_token is in database
	user := structs.User{}
	user.Session = structs.Session{AccessToken: accessToken}

	globals.DBConn.Preload("Session").First(&user, user)
	if user.Session.ID == 0 {
		// Redirect to login
		return c.Redirect("/oauth")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(user.RippleId),
	})
	return nil
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork: false, // true = multithreaded, false = singlethreaded
	})

	app.Static("/", "/frontend")
	app.Static("/admin", "/admin")

	app.Post("/api/v1/upload", routes.Upload)
	app.Get("/oauth", routes.GetOauth)

	// use auth middleware
	app.Use(checkSessionCookie)

	app.Listen(":5000")
}
