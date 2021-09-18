package routes

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func GetOAuthRipple(c *fiber.Ctx) error {
	sess, err := globals.SessionStore.Get(c)
	if err != nil {
		return err
	}

	// Get the oauth2 config
	oauthConfig := oauth2.Config{
		ClientID:     os.Getenv("RIPPLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("RIPPLE_OAUTH_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://ripple.moe/oauth/authorize",
			TokenURL:  "https://ripple.moe/oauth/token",
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		RedirectURL: os.Getenv("RIPPLE_OAUTH_REDIRECT_URL"),
		Scopes:      []string{""},
	}

	// Get the code
	code := c.Query("code")

	if code == "" {
		id := sess.ID()
		sess.Set("oauth_state", id)
		// Save session
		if err := sess.Save(); err != nil {
			return err
		}
		globals.Logger.WithFields(logrus.Fields{
			"get":         sess.Get(id),
			"oauth_state": id,
		}).Info("Saved session")

		return c.Redirect(oauthConfig.AuthCodeURL(id))
	}

	token, err := getOauth(c, &oauthConfig, code)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get token",
			"error":   err.Error(),
		})
	}

	rippleResp := structs.RippleSelf{}

	globals.Logger.WithField("token", token).Debug("Attempting to receive /api/v1/users/self")
	client := oauthConfig.Client(context.Background(), token)

	// Get the user info
	resp, err := client.Get("https://ripple.moe/api/v1/users/self")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get user data",
			"error":   err.Error(),
		})
	}

	json.NewDecoder(resp.Body).Decode(&rippleResp)

	// check if rippleResp has the UserId
	if rippleResp.UserId == 0 {
		globals.Logger.WithField("token", token).Debug("Failed getting /api/v1/users/self")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success":         false,
			"message":         "Could not get user data after trying 3 times",
			"error":           "UserId is empty",
			"responseDecoded": rippleResp,
			"token":           token,
		})
	}

	user := structs.NewUser()

	// Check for session token
	sessionToken, _ := globals.CheckAuth(c)
	search := structs.NewUser()
	search.Session.SessionToken = sessionToken
	res := structs.NewUser()

	globals.DBConn.Debug().Preload("Session").First(&res, search)

	// check if a user was found
	if res.ID != 0 {
		user = res
	}

	user.Session.SessionToken = sess.ID()

	globals.DBConn.Debug().Save(&user)

	// Save session
	if err := sess.Save(); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).Redirect("/")
}

func getOauth(c *fiber.Ctx, oauthConfig *oauth2.Config, code string) (*oauth2.Token, error) {
	sess, err := globals.SessionStore.Get(c)
	if err != nil {
		return nil, err
	}

	globals.Logger.WithFields(logrus.Fields{
		"code": code,
	}).Debug("Received code")

	// Read oauthState from Cookie
	oauth_state := sess.Get("oauth_state")

	if c.Query("state") != oauth_state {
		globals.Logger.WithFields(logrus.Fields{
			"state":       c.Query("state"),
			"oauth_state": oauth_state,
			"connection":  c.IP(),
		}).Error("State mismatch")

		return nil, fiber.ErrBadRequest
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return token, nil
}
