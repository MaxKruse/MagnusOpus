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

	err := localDB.Joins("JOIN staffs ON staffs.tournament_id = tournaments.id JOIN users ON users.id = staffs.user_id").Preload("Staffs").Preload("Rounds").Where("visible = ?", true).Or("users.ripple_id = ?", self.RippleId).First(&tournament, c.Params("id")).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	return c.Status(fiber.StatusOK).JSON(tournament)
}
