package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetTournament(c *fiber.Ctx) error {
	self, err := utils.GetSelfFromSess(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// check if we are a staff member
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	tournament, err := utils.GetTournament(uint(id))

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotFound)
	}

	_, err = utils.CanViewTournament(self.ID, tournament.ID)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusForbidden)
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
