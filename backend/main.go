package main

import (
	// Import logrus

	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	psql "github.com/gofiber/storage/postgres"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/routes"
	"github.com/maxkruse/magnusopus/backend/routes/tournaments"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

func init() {
	defer utils.TimeTrack(time.Now(), "init")
	// Setup global state before main

	// Create a new instance of the logger
	globals.Logger = logrus.New()

	// Set the log level
	globals.Logger.SetLevel(logrus.DebugLevel)

	// Set log config
	globals.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.1234",
		PrettyPrint:     true,
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
	localDB := globals.DBConn

	// Migrate Tables
	localDB.AutoMigrate(&structs.User{})
	localDB.AutoMigrate(&structs.Session{})
	localDB.AutoMigrate(&structs.Tournament{})
	localDB.AutoMigrate(&structs.Staff{})
	localDB.AutoMigrate(&structs.Round{})
	globals.Logger.Debug("Migrated")

	globals.SessionStore = session.New(session.Config{
		Expiration: time.Hour * 24 * 7 * 31, // 1 Month
		Storage: psql.New(psql.Config{
			Reset:    false,
			Host:     "session_store",
			Port:     5432,
			Username: globals.Config.POSTGRES_USER,
			Password: globals.Config.POSTGRES_PASSWORD,
			Database: globals.Config.POSTGRES_DB,
			Table:    "sessions",
		}),
	})

	// Fill in the Ripple UserId for superadmins. Only those can create tournaments
	globals.AllowedSuperadmin = []int{1955}

	globals.Logger.Info("Starting Magnusopus backend")
	globals.Logger.WithFields(logrus.Fields{"config": globals.Config}).Debug("Config")
}

func checkSessionCookie(c *fiber.Ctx) error {
	defer utils.TimeTrack(time.Now(), "checkSessionCookie")

	// dont apply to /api/v1/tournaments
	if c.Path() == "/api/v1/tournaments" && c.Method() == "GET" {
		return c.Next()
	}

	user, err := utils.GetSelf(c)

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	globals.Logger.WithField("user", user).Debug("User")

	return c.Next()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Magnusopus Backend",
	})

	// use middlewares
	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05.1234",
		TimeZone:   "UTC",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// oauth routes
	oauth := app.Group("/oauth")
	oauth.Get("/ripple", routes.GetOAuthRipple)

	api := app.Group("/api")

	// Status Page
	app.Get("/status", monitor.New())

	// use custom middleware for the entire api
	api.Use(checkSessionCookie)

	// add api v1 (backwards compatability stuff)
	v1 := api.Group("/v1")

	v1.Get("/me", routes.Me)
	v1.Get("/self", routes.Me)

	// Tournament
	v1.Get("/tournaments", tournaments.GetTournaments)
	v1.Get("/tournaments/:id", tournaments.GetTournament)

	v1.Post("/tournaments", tournaments.PostTournament)

	globals.Logger.Fatal(app.Listen(":5000"))
}
