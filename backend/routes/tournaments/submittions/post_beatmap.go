package submittions

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func Upload(c *fiber.Ctx) error {
	localDB := globals.DBConn
	self, err := utils.GetSelfFromDB(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// check if user is registered to tournament
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	tournament, err := utils.GetTournament(uint(id))

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotFound)
	}

	err = tournament.IsRegistered(localDB, self.ID)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusForbidden)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnprocessableEntity)
	}

	// get md5 of file
	md5Gen := sha256.New()
	md5Gen.Write([]byte(time.Now().String() + fmt.Sprint(self.ID)))

	md5FileName := md5Gen.Sum(nil)[:16]
	fileName := hex.EncodeToString(md5FileName)

	filePath := fmt.Sprintf("/storage/%s.osu", fileName)
	// check if filePath exists
	if _, err := os.Stat(filePath); err == nil {
		return utils.DefaultErrorMessage(c, errors.New("file already exists"), fiber.StatusUnprocessableEntity)
	}

	// TODO(Lithium): Check if file is allowed

	// read first line of file
	f, err := file.Open()
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()              // Scan the first line
	firstLine := scanner.Text() // first line of file

	if !strings.HasPrefix(firstLine, "osu file format v") {
		return utils.DefaultErrorMessage(c, errors.New("file is not a valid osu file"), fiber.StatusUnprocessableEntity)
	}

	// TODO(Lithium): Check if user has too many files
	var count int64
	localDB.Model(&structs.BeatmapSubmittion{}).Where("user_id = ?", self.ID).Count(&count)

	if count >= globals.MAX_FILES_PER_USER {
		return utils.DefaultErrorMessage(c, errors.New("you have too many files, delete at least one"), fiber.StatusUnprocessableEntity)
	}

	// TODO(Lithium): Save BeatmapSubmittion for this user and round
	err = c.SaveFile(file, filePath)
	if err != nil {
		os.Remove(filePath)
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	round := structs.Round{}
	err = localDB.Where("name = ? AND tournament_id = ?", c.Params("name"), c.Params("id")).First(&round).Error
	if err != nil {
		os.Remove(filePath)
		return utils.DefaultErrorMessage(c, errors.New("no round found"), fiber.StatusInternalServerError)
	}

	submittion := structs.BeatmapSubmittion{
		Round:        &round,
		User:         &self,
		Hash:         fileName,
		DownloadPath: fmt.Sprintf("/download/%s.osu", fileName),
		ToUse:        true,
	}

	err = localDB.Create(&submittion).Error
	if err != nil {
		os.Remove(filePath)
		return utils.DefaultErrorMessage(c, errors.New("cannot save submittion, "+err.Error()), fiber.StatusInternalServerError)
	}

	// omit fields for response
	submittion.Round = nil
	submittion.User = nil

	return c.Status(fiber.StatusCreated).JSON(submittion)
}
