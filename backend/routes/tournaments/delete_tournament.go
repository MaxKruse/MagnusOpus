package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func DeleteTournament(c *fiber.Ctx) error {
	localDB := globals.DBConn
	id := c.Params("id", "")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	self, err := utils.GetSelfFromSess(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	if err := utils.CanEditTournament(self.ID, uint(idUint)); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnauthorized)
	}

	t := structs.Tournament{}
	err = localDB.Delete(&t, id).Error

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "deleted tournament",
	})

}
