package submittions

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func Upload(c *fiber.Ctx) error {
	sess, err := globals.SessionStore.Get(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnprocessableEntity)
	}

	// get md5 of file
	md5Gen := sha256.New()
	md5Gen.Write([]byte(time.Now().String() + sess.ID()))

	md5FileName := md5Gen.Sum(nil)[:16]
	fileName := hex.EncodeToString(md5FileName)
	log.Println(fileName)

	filePath := fmt.Sprintf("/storage/%s.osu", fileName)
	// check if filePath exists
	if _, err := os.Stat(filePath); err == nil {
		return utils.DefaultErrorMessage(c, errors.New("file already exists"), fiber.StatusUnprocessableEntity)
	}

	// TODO(Lithium): Check if file is allowed
	// TODO(Lithium): Check if user has too many files
	// TODO(Lithium): Save BeatmapSubmittion for this user and round
	err = c.SaveFile(file, filePath)
	if err != nil {
		os.Remove(filePath)
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
	})
}
