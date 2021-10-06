package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
)

func GetUsers(c *fiber.Ctx) error {
	localDB := globals.DBConn

	var users []structs.User
	localDB.Find(&users)
	return c.Status(fiber.StatusOK).JSON(users)
}
