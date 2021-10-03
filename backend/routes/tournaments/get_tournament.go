package tournaments

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetTournament(c *fiber.Ctx) error {
	tournament := structs.Tournament{}

	self, _ := utils.GetSelf(c)

	// check if we are a staff member
	log.Println(c.Params("id"))
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	log.Println(id)
	log.Println(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid tournament id",
			"success": false,
		})
	}

	tournament, err = utils.GetTournament(uint(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	canView := utils.CanViewTournament(self.ID, tournament.ID)

	if canView != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   canView.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
