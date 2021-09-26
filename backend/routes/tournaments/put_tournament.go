package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func PutTournament(c *fiber.Ctx) error {
	if !utils.IsSuperadmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"success": false,
		})
	}
	c.Accepts("application/json")
	t := structs.Tournament{}

	// Decode body
	err := c.BodyParser(&t)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// TODO: Validate tournament
	err = validTournament(t)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// get tournament from id in db
	res := structs.Tournament{}
	err = globals.DBConn.First(&res, c.Params("id")).Error

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

	localDB := globals.DBConn
	err = localDB.Where(update.ID).UpdateColumns(&t).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(update)
}
