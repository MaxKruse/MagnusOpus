package routes

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"golang.org/x/oauth2"
)

func state(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func GetOauth(c *fiber.Ctx) error {
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
		// Redirect to the oauth page
		state, err := state(32)
		if err != nil {
			return err
		}

		// make new cookie with state value
		cookie := &fiber.Cookie{
			Name:  "oauth_state",
			Value: state,
		}
		c.Cookie(cookie)
		return c.Redirect(oauthConfig.AuthCodeURL(state))
	}

	// Read oauthState from Cookie
	oauth_state := c.Cookies("oauth_state")

	if c.Query("state") != oauth_state {
		log.Println("invalid oauth state")
		return c.Status(fiber.StatusUnauthorized).SendString("invalid oauth state")
		// return c.Redirect("/")
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get token",
			"error":   err.Error(),
		})
	}

	// Get the user info
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://ripple.moe/api/v1/ping")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get user info",
			"error":   err.Error(),
		})
	}

	defer resp.Body.Close()

	rippleResp := structs.RipplePing{}
	json.NewDecoder(resp.Body).Decode(&rippleResp)

	user := structs.User{}
	user.RippleId = rippleResp.UserId

	// check if user with this RippleId exists
	globals.DBConn.Preload("Session").First(&user, user)

	// check if user had a session
	if user.Session.ID != 0 {
		// update the session
		user.Session.AccessToken = token.AccessToken
		err = globals.DBConn.Save(&user.Session).Error

	} else {
		user.Session = structs.Session{
			AccessToken: token.AccessToken,
		}
		err = globals.DBConn.Save(&user).Error

	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not save user",
			"error":   err.Error(),
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not unmarshal user info",
			"error":   err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "ripple_token",
		Value: user.Session.AccessToken,
	})

	// Clean up after ourselfs
	c.ClearCookie("oauth_state")

	return c.Status(fiber.StatusOK).Redirect("/")
}
