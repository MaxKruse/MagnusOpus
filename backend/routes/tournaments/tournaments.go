package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetTournaments(c *fiber.Ctx) error {
	tournaments := []structs.Tournament{}

	globals.DBConn.Preload("Round").Find(&tournaments)

	return c.Status(fiber.StatusOK).JSON(tournaments)
}

func GetTournament(c *fiber.Ctx) error {
	tournament := structs.Tournament{}

	globals.DBConn.Preload("User").Preload("Staff").Preload("Round").First(&tournament, c.Params("id"))

	return c.Status(fiber.StatusOK).JSON(tournament)
}
