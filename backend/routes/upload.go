package routes

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func Upload(c *fiber.Ctx) error {
	// check if authenticated
	if _, err := utils.CheckAuth(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"file":    file,
		})
	}

	filePath := fmt.Sprintf("/storage/%s.osu", file.Filename)
	// check if filePath exists
	if _, err := os.Stat(filePath); err == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "File already exists",
		})
	}

	// TODO(Lithium): Check if file is allowed
	// TODO(Lithium): Check if user has too many files
	err = c.SaveFile(file, filePath)
	if err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
	})
}
