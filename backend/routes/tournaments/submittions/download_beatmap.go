package submittions

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func DownloadMap(c *fiber.Ctx) error {
	fileName := c.Params("file")

	// check if file exists in /storage
	filePath := fmt.Sprintf("/storage/%s.osu", fileName)
	if _, err := os.Stat(filePath); err != nil {
		return utils.DefaultErrorMessage(c, errors.New("file not found"), fiber.StatusNotFound)
	}

	// get BeatmapSubmittion from filename
	var beatmap structs.BeatmapSubmittion
	downloadPath := fmt.Sprintf("/download/%s.osu", fileName)

	localDB := globals.DBConn
	err := localDB.Preload("Round").Preload("User").Where("download_path = ?", downloadPath).First(&beatmap).Error
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	self, _ := utils.GetSelfFromDB(c)
	downloadName := fmt.Sprintf("%s - %s.osu", beatmap.Round.Name, beatmap.Hash)

	// check if we submitted it
	if beatmap.User.ID == self.ID {
		return c.Download(filePath, downloadName)
	}

	// check the associated tournament from the round and see if we are a judge or above
	var tournament structs.Tournament
	err = localDB.Where("id = ?", beatmap.Round.ID).First(&tournament).Error
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// check if we are a judge or above
	if err := utils.CanJudge(self.ID, tournament.ID); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusForbidden)
	}

	return c.Download(filePath, downloadName)
}
