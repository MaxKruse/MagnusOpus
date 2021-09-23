package tournaments

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func validTournament(t structs.Tournament) error {

	if t.Name == "" {
		return errors.New("name is required")
	}

	if t.Description == "" {
		return errors.New("description is required")
	}

	if t.DownloadPath == "" {
		return errors.New("download_path is required")
	}

	zeroTime := time.Time{}

	if t.StartTime == zeroTime {
		return errors.New("start_time is required")
	}

	if t.EndTime == zeroTime {
		return errors.New("end_time is required")
	}

	return nil
}

func PostTournament(c *fiber.Ctx) error {
	defer utils.TimeTrack(time.Now(), "PostTournament")
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

	// TODO: Validate tournament
	err = validTournament(t)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// Dont Accept: Rounds, Staff. Instead, set them to nil
	t.Rounds = nil
	t.Staffs = nil

	// TODO: Check if tournament exists
	localDB := globals.DBConn
	err = localDB.Where("name = ?", t.Name).Error

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	// TODO: Add current user as tournament staff and owner

	me, err := utils.GetSelf(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}
	staffMember := structs.Staff{
		UserId: me.ID,
		Role:   "owner",
	}

	t.Staffs = append(t.Staffs, staffMember)

	// TODO: Save tournament
	err = localDB.Save(&t).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(t)
}
