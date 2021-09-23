package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/structs"
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
