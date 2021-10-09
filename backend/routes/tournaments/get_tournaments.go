package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetTournaments(c *fiber.Ctx) error {
	tournaments := []*structs.Tournament{}
	results := tournaments
	localDB := globals.DBConn
	self, err := utils.GetSelfFromSess(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	localDB.Preload("Staffs").Find(&tournaments)

	for _, tournament := range tournaments {
		canView := tournament.Visible
		for _, staff := range tournament.Staffs {
			if staff.UserId == self.ID {
				canView = true
			}
		}

		if canView {
			tournament.Staffs = nil
			results = append(results, tournament)
		}
	}

	return c.Status(fiber.StatusOK).JSON(results)
}
