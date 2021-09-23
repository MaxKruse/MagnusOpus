package tournaments

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetTournaments(c *fiber.Ctx) error {
	defer utils.TimeTrack(time.Now(), "GetTournaments")
	tournaments := []*structs.Tournament{}
	filter := utils.GetRequestFilter(c)
	localDB := globals.DBConn

	localDB.Preload("Staffs").Limit(filter.Limit).Offset(filter.Offset).Find(&tournaments)

	return c.Status(fiber.StatusOK).JSON(tournaments)
}
