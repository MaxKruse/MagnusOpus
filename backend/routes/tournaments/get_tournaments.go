package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/util"
)

func GetTournaments(c *fiber.Ctx) error {
	tournaments := []*structs.Tournament{}
	filter := util.GetRequestFilter(c)
	localDB := globals.DBConn

	localDB.Preload("Staffs").Limit(filter.Limit).Offset(filter.Offset).Find(&tournaments)

	return c.Status(fiber.StatusOK).JSON(tournaments)
}
