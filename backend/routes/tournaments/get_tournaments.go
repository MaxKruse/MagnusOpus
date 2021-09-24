package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
	"github.com/sirupsen/logrus"
)

func GetTournaments(c *fiber.Ctx) error {
	visible_tournaments := []*structs.Tournament{}
	staff_tournaments := []*structs.Tournament{}
	localDB := globals.DBConn
	self, _ := utils.GetSelf(c)

	localDB.Where("visible = ?", true).Find(&visible_tournaments)
	localDB.Joins("LEFT JOIN staffs ON staffs.tournament_id = tournaments.id JOIN users ON users.id = staffs.user_id").Where("users.ripple_id = ?", self.RippleId).Find(&staff_tournaments)

	globals.Logger.WithFields(logrus.Fields{
		"visible_tournaments": visible_tournaments,
		"staff_tournaments":   staff_tournaments,
	}).Info("GetTournaments")

	return c.Status(fiber.StatusOK).JSON(append(visible_tournaments, staff_tournaments...))
}
