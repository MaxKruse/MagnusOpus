package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func DeleteTournament(c *fiber.Ctx) error {
	localDB := globals.DBConn
	id := c.Params("id", "")

	t := structs.Tournament{}
	err := localDB.Delete(&t, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "deleted tournament",
	})

}
