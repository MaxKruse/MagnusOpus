package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetTournament(c *fiber.Ctx) error {
	tournament := structs.Tournament{}

	err := globals.DBConn.Preload("User").Preload("Staff").Preload("Round").First(&tournament, c.Params("id")).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
