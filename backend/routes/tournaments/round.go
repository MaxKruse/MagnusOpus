package tournaments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
)

func ActivateRound(c *fiber.Ctx) error {
	self, err := utils.GetSelfFromDB(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	tournamentID := c.Params("id")
	tournamentIDUint, err := strconv.ParseUint(tournamentID, 10, 64)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	editErr := utils.CanEditRounds(self.ID, uint(tournamentIDUint))
	if editErr != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusUnauthorized)
	}

	round := structs.Round{}
	if err := c.BodyParser(&round); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	t, err := utils.GetTournament(uint(tournamentIDUint))
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusNotFound)
	}

	err = t.ActivateRound(round.Name)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusBadRequest)
	}

	localDB := globals.DBConn
	localDB.Save(&t)

	return c.Status(fiber.StatusCreated).JSON(t)
}

func AddRound(c *fiber.Ctx) error {
	self, err := utils.GetSelfFromDB(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	tournamentID := c.Params("id")
	tournamentIDUint, err := strconv.ParseUint(tournamentID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid tournament ID",
			"success": false,
		})
	}

	editErr := utils.CanEditRounds(self.ID, uint(tournamentIDUint))
	if editErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   editErr.Error(),
			"success": false,
		})
	}

	round := structs.Round{}
	if err := c.BodyParser(&round); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	localDB := globals.DBConn
	err = round.ValidRound(localDB)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	err = round.RoundExist(localDB)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Round already exists",
			"success": false,
		})
	}

	t, err := utils.GetTournament(uint(tournamentIDUint))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	t.Rounds = append(t.Rounds, round)
	// Activate this round
	err = t.ActivateRound(round.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	localDB.Save(&t)

	return c.Status(fiber.StatusOK).JSON(t)
}
