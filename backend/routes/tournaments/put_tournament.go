package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func PutTournament(c *fiber.Ctx) error {
	self, err := utils.GetSelfFromSess(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	tournament_id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	tournament_id := uint(tournament_id64)

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	if editErr := utils.CanEditTournament(self.ID, tournament_id); editErr != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnauthorized)
	}
	c.Accepts("application/json")
	t := structs.Tournament{}

	// Decode body
	err = c.BodyParser(&t)

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	localDB := globals.DBConn
	err = t.ValidTournament(localDB)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnprocessableEntity)
	}

	// get tournament from id in db
	res := structs.Tournament{}
	err = localDB.First(&res, tournament_id).Error

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotFound)
	}

	// Dont Accept: Rounds, Staff. Instead, set them to nil
	t.Rounds = nil
	t.Staffs = nil

	// update tournament
	update := t
	update.ID = res.ID
	update.Rounds = res.Rounds
	update.Staffs = res.Staffs

	err = localDB.Where(update.ID).UpdateColumns(&t).Error

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(update)
}
