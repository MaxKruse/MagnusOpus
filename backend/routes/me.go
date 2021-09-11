package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
)

func Me(c *fiber.Ctx) error {
	self, err := globals.GetSelf(c)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(self)
}
