package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func Register(c *fiber.Ctx) error {
	self, err := utils.GetSelf(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err)
	}

	// get current tournament
	id := c.Params("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return utils.DefaultErrorMessage(c, err)
	}

	tournament, err := utils.GetTournament(uint(idUint))
	if err != nil {
		return utils.DefaultErrorMessage(c, err)
	}

	tournament.Registrations = append(tournament.Registrations, self)

	localDB := globals.DBConn
	err = localDB.Save(&tournament).Error
	if err != nil {
		return utils.DefaultErrorMessage(c, err)
	}

	// nil staff and rounds before returning
	tournament.Staffs = nil
	tournament.Rounds = nil

	return c.Status(fiber.StatusCreated).JSON(tournament)
}
