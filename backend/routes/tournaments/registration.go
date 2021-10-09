package tournaments

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func Register(c *fiber.Ctx) error {
	self, err := utils.GetSelfFromDB(c)
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

	// check if this tournament is even available
	if _, err := utils.CanViewTournament(self.ID, tournament.ID); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusForbidden)
	}

	if err := tournament.IsRegistered(localDB, self.ID); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotAcceptable)
	}

	if err := tournament.RegistrationsOpen(); err != nil {
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

func Unregister(c *fiber.Ctx) error {
	self, err := utils.GetSelfFromSess(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnauthorized)
	}

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

	localDB := globals.DBConn

	// check if we are actually registered
	if err := tournament.IsRegistered(localDB, self.ID); err == nil {
		return utils.DefaultErrorMessage(c, errors.New("not registered to tournament"), fiber.StatusNotAcceptable)
	}

	// check if registrations are still open
	if err := tournament.RegistrationsOpen(); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotAcceptable)
	}

	// remove self from registrations
	for _, registration := range tournament.Registrations {
		if registration.ID == self.ID {
			localDB.Model(&tournament).Association("Registrations").Delete(&registration)
			break
		}
	}

	tournament, _ = utils.GetTournament(tournament.ID)

	return c.Status(fiber.StatusOK).JSON(tournament)
}
