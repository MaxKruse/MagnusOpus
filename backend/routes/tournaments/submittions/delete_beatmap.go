package submittions

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func DeleteBeatmap(c *fiber.Ctx) error {
	localDB := globals.DBConn

	deleteStr := c.Params("to_delete", "")
	deleteID, err := utils.StringToUint32(deleteStr)
	if err != nil {
		return utils.DefaultErrorMessage(c, fiber.ErrBadRequest, fiber.StatusBadRequest)
	}

	if deleteID == 0 {
		return utils.DefaultErrorMessage(c, fiber.ErrBadRequest, fiber.StatusBadRequest)
	}

	self, err := utils.GetSelfFromDB(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// check if we own this to_delete
	submittion, err := self.OwnsMap(localDB, deleteID)
	if err != nil {
		return utils.DefaultErrorMessage(c, fiber.ErrUnauthorized, fiber.StatusUnauthorized)
	}

	err = localDB.Delete(&submittion, submittion.ID).Error
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	err = os.Remove(submittion.DownloadPath)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "deleted submittion",
	})

}
