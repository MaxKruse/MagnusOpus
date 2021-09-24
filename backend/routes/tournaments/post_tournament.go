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
	defer utils.TimeTrack(time.Now(), "validTournament")

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
		return errors.New("start_time is required (ISO 8601) (RFC 3339)")
	}

	if t.EndTime == zeroTime {
		return errors.New("end_time is required (ISO 8601) (RFC 3339)")
	}

	// Check if the time is in the future
	if t.StartTime.Before(time.Now()) {
		return errors.New("start_time must be in the future")
	}

	if t.EndTime.Before(time.Now()) {
		return errors.New("end_time must be in the future")
	}

	// check if end_time is at least 3 days after start_time
	if t.EndTime.Sub(t.StartTime) < (3 * 24 * time.Hour) {
		return errors.New("end_time must be at least 3 days after start_time")
	}

	localDB := globals.DBConn
	res := structs.Tournament{}
	localDB.Where(t).First(&res)

	if res.ID != 0 {
		return errors.New("name must be unique")
	}

	return nil
}

func PostTournament(c *fiber.Ctx) error {
	if !utils.IsSuperadmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"success": false,
		})
	}
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

	localDB := globals.DBConn
	err = localDB.Save(&t).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(t)
}
