package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetTournament(c *fiber.Ctx) error {
	selfID, err := utils.GetSelfID(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// check if we are a staff member
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid tournament id",
			"success": false,
		})
	}

	tournament, err := utils.GetTournament(uint(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	if canView := utils.CanViewTournament(selfID, tournament.ID); canView != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   canView.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
