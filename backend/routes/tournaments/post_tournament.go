package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func PostTournament(c *fiber.Ctx) error {
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

	localDB := globals.DBConn

	// TODO: Validate tournament
	err = t.ValidTournament(localDB)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// Dont Accept: Rounds, Staff. Instead, set them to nil
	t.Rounds = nil
	t.Staffs = nil

	me, err := utils.GetSelf(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	staffMember := structs.Staff{
		User: &me,
		Role: "owner",
	}
	t.Staffs = append(t.Staffs, staffMember)

	err = localDB.Save(&t).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// Remove session from staff user to not display
	for _, staff := range t.Staffs {
		staff.User.Sessions = nil
	}

	return c.Status(fiber.StatusOK).JSON(t)
}
