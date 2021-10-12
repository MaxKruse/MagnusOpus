package tournaments

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func PostTournamentStaff(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	id, err := utils.StringToUint32(tournamentID)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	self, err := utils.GetSelfFromDB(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	adminErr := utils.CanAddStaff(self.ID, uint(id))
	if adminErr != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnauthorized)
	}

	localDB := globals.DBConn

	staffReq := structs.StaffPost{}
	if err := c.BodyParser(&staffReq); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	if err := structs.ValidStaff(staffReq); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	t, err := utils.GetTournament(uint(id))
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotFound)
	}

	user := structs.User{}
	err = localDB.First(&user, staffReq.UserId).Error
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotFound)
	}

	staff := structs.Staff{
		User: &user,
		Role: staffReq.Role,
	}
	t.Staffs = append(t.Staffs, staff)

	// Save tournament
	err = localDB.Save(&t).Error
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(t.Staffs)
}
