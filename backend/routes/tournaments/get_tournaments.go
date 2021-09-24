package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetTournaments(c *fiber.Ctx) error {
	tournaments := []*structs.Tournament{}
	localDB := globals.DBConn
	self, _ := utils.GetSelf(c)

	localDB.Joins("JOIN staffs ON staffs.tournament_id = tournaments.id JOIN users ON users.id = staffs.user_id").Preload("Rounds").Where("visible = ?", true).Or("users.ripple_id = ?", self.RippleId).Find(&tournaments)

	return c.Status(fiber.StatusOK).JSON(tournaments)
}
