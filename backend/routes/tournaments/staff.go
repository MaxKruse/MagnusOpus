package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func PostTournamentStaff(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	id, err := strconv.ParseUint(tournamentID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"sucess": false,
			"error":  err.Error(),
		})
	}

	selfID, err := utils.GetSelfID(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	adminErr := utils.CanAddStaff(selfID, uint(id))
	if adminErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sucess": false,
			"error":  adminErr.Error(),
		})
	}

	localDB := globals.DBConn

	staffReq := structs.StaffPost{}
	if err := c.BodyParser(&staffReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"sucess": false,
			"error":  err.Error(),
		})
	}

	if err := structs.ValidStaff(staffReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"sucess": false,
			"error":  err.Error(),
		})
	}

	t, err := utils.GetTournament(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"sucess": false,
			"error":  err.Error(),
		})
	}

	user := structs.User{}
	err = localDB.First(&user, staffReq.UserId).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"sucess": false,
			"error":  err.Error(),
		})
	}

	staff := structs.Staff{
		User: &user,
		Role: staffReq.Role,
	}
	t.Staffs = append(t.Staffs, staff)

	// Save tournament
	err = localDB.Save(&t).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"sucess": false,
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(t.Staffs)
}
