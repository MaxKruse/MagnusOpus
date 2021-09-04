package main

import (
	// Import logrus

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/routes"
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
	}

	globals.Logger.Debug("Connected to database")

	globals.Logger.Debug("Starting Magnusopus backend")
	globals.Logger.WithFields(logrus.Fields{"config": globals.Config}).Debug("Config")
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork: false, // true = multithreaded, false = singlethreaded
	})

	app.Post("/upload", routes.Upload)

	app.Listen(":5000")
}
