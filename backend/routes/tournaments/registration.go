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
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnauthorized)
	}

	// nil some self values for response
	self.Sessions = nil

	// get current tournament
	id := c.Params("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	tournament, err := utils.GetTournament(uint(idUint))
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// check if we are already registered
	localDB := globals.DBConn

	if err := tournament.IsRegistered(localDB, self.ID); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotAcceptable)
	}

	tournament.Registrations = append(tournament.Registrations, self)
	err = localDB.Save(&tournament).Error
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// nil staff and rounds before returning
	tournament.Staffs = nil
	tournament.Rounds = nil

	return c.Status(fiber.StatusCreated).JSON(tournament)
}
