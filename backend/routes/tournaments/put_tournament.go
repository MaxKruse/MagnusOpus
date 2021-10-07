package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func PutTournament(c *fiber.Ctx) error {
	selfID, err := utils.GetSelfID(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	tournament_id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	tournament_id := uint(tournament_id64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "Invalid tournament id",
			"success": false,
		})
	}

	if editErr := utils.CanEditTournament(selfID, tournament_id); editErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   editErr.Error(),
			"success": false,
		})
	}
	c.Accepts("application/json")
	t := structs.Tournament{}

	// Decode body
	err = c.BodyParser(&t)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	localDB := globals.DBConn
	err = t.ValidTournament(localDB)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// get tournament from id in db
	res := structs.Tournament{}
	err = localDB.First(&res, tournament_id).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// Dont Accept: Rounds, Staff. Instead, set them to nil
	t.Rounds = nil
	t.Staffs = nil

	// update tournament
	update := t
	update.ID = res.ID
	update.Rounds = res.Rounds
	update.Staffs = res.Staffs

	globals.Logger.WithField("tournament", update).Info("Updating tournament")

	err = localDB.Where(update.ID).UpdateColumns(&t).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(update)
}
