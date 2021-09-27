package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetTournament(c *fiber.Ctx) error {
	tournament := structs.Tournament{}
	localDB := globals.DBConn

	self, _ := utils.GetSelf(c)

	// check if we are a staff member
	err := localDB.Preload("Staffs.User").Preload("Rounds").First(&tournament, c.Params("id")).Error

	canView := tournament.Visible
	for _, staff := range tournament.Staffs {
		if staff.UserId == self.ID {
			canView = true
		}
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	if !canView {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
