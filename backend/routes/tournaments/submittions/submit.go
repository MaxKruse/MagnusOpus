package submittions

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func GetBeatmaps(c *fiber.Ctx) error {
	// NOTE:
	// This path can have multiple paths:
	// - Just a user trying to get their own maps
	// - A Judge trying to get maps for a round

	self, err := utils.GetSelf(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnauthorized)
	}

	// get current tournament
	id := c.Params("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	tournament, err := utils.GetTournament(uint(idUint))
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// check if this tournament is even available
	if _, err := utils.CanViewTournament(self.ID, tournament.ID); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusForbidden)
	}

	maps := []structs.BeatmapSubmittion{}
	err = nil
	localDB := globals.DBConn

	roundName := c.Params("name")

	// branch paths as appropriate
	if tmpErr := utils.CanJudge(self.ID, tournament.ID); tmpErr == nil {
		maps, err = tournament.GetBeatmapsForJudge(localDB, roundName)
	} else {
		maps, err = tournament.GetBeatmapsForUser(localDB, self.ID, roundName)
	}

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(maps)
}
