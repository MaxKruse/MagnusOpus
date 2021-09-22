package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/util"
)

func PostTournament(c *fiber.Ctx) error {
	c.Accepts("application/json")
	t := structs.Tournament{}
	c.BodyParser(&t)

	// TODO: Validate tournament
	// TODO: Check if tournament exists
	// TODO: Save tournament
	// TODO: Add current user as tournament staff and owner

	return c.Status(fiber.StatusOK).JSON(t)
}

func GetTournaments(c *fiber.Ctx) error {
	tournaments := []*structs.Tournament{}
	filter := util.GetRequestFilter(c)

	globals.DBConn.Preload("Staffs").Limit(filter.Limit).Offset(filter.Offset).Find(&tournaments)

	return c.Status(fiber.StatusOK).JSON(tournaments)
}

func GetTournament(c *fiber.Ctx) error {
	tournament := structs.Tournament{}

	err := globals.DBConn.Debug().Preload("User").Preload("Staff").Preload("Round").First(&tournament, c.Params("id")).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
