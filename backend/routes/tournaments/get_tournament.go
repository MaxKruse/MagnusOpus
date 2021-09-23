package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetTournament(c *fiber.Ctx) error {
	tournament := structs.Tournament{}
	localDB := globals.DBConn

	err := localDB.Preload("Staffs").Preload("Rounds").First(&tournament, c.Params("id")).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
