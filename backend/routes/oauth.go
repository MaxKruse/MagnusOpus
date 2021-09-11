package routes

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		globals.Logger.Fatalf("crypto/rand is unavailable: Read() failed with %#v", err)
	}
}

func generateRandomString(n int) ([]byte, error) {
	assertAvailablePRNG()

	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func genSessionToken(n int) (string, error) {
	data, err := generateRandomString(n)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func GetOAuthRipple(c *fiber.Ctx) error {
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
		state, err := genSessionToken(32)
		if err != nil {
			return err
		}

		// make new cookie with state value
		cookie := &fiber.Cookie{
			Name:  "oauth_state",
			Value: state,
		}
		c.Cookie(cookie)

		globals.Logger.WithFields(logrus.Fields{
			"state":    state,
			"redirect": oauthConfig.AuthCodeURL(state),
		}).Debug("Generated new state, redirecting")

		return c.Redirect(oauthConfig.AuthCodeURL(state))
	}

	token, err := getOauth(c, &oauthConfig, code)

	globals.Logger.WithFields(logrus.Fields{
		"token": token,
	}).Info("Got token")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get token",
			"error":   err.Error(),
		})
	}

	// Get the user info
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://ripple.moe/api/v1/users/self")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get user data",
			"error":   err.Error(),
		})
	}

	defer resp.Body.Close()

	rippleResp := structs.RippleSelf{}
	json.NewDecoder(resp.Body).Decode(&rippleResp)

	globals.Logger.WithFields(logrus.Fields{
		"user": rippleResp,
	}).Info("Got user data")

	user := structs.User{}
	user.RippleId = rippleResp.UserId

	// Check for session token
	sessionToken, _ := globals.CheckAuth(c)
	search := structs.User{
		Session: &structs.Session{
			SessionToken: sessionToken,
		},
	}
	res := structs.User{}

	globals.DBConn.Preload("Session").First(&res, search)

	// Associate ripple accounts if a previous session existed for a bancho user
	if res.RippleId != 0 {
		globals.Logger.WithFields(logrus.Fields{
			"res": res,
		}).Debug("Searched for session and found user")
		user = res
		user.RippleId = rippleResp.UserId
		user.Username = rippleResp.Username

		user.Session.AccessToken = token.AccessToken
		user.Session.RefreshToken = token.RefreshToken

		globals.Logger.WithFields(logrus.Fields{
			"user": user,
		}).Debug("Saving user")

		err = globals.DBConn.Save(&user).Error

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Could not save user",
				"error":   err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:  "session_token",
			Value: user.Session.SessionToken,
		})

		// Clean up after ourselfs
		c.ClearCookie("oauth_state")

		return c.Status(fiber.StatusOK).Redirect("/")

	} else {
		sessionToken, err = genSessionToken(32)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Could not generate session token",
				"error":   err.Error(),
			})
		}
	}

	// check if user with this BanchoId exists
	globals.DBConn.Preload("Session").First(&res, structs.User{BanchoId: user.BanchoId})

	// check if user had a session
	if res.Session.ID != 0 {
		// update the session
		user.Session.SessionToken = sessionToken
		user.Session.AccessToken = token.AccessToken
		user.Session.RefreshToken = token.RefreshToken
		err = globals.DBConn.Save(&user.Session).Error

	} else {
		user.Session = &structs.Session{
			AccessToken:  token.AccessToken,
			SessionToken: sessionToken,
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
		Name:  "session_token",
		Value: user.Session.SessionToken,
	})

	// Clean up after ourselfs
	c.ClearCookie("oauth_state")

	return c.Status(fiber.StatusOK).Redirect("/")
}

func GetOAuthBancho(c *fiber.Ctx) error {
	// Get the oauth2 config
	oauthConfig := oauth2.Config{
		ClientID:     os.Getenv("BANCHO_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("BANCHO_OAUTH_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://osu.ppy.sh/oauth/authorize",
			TokenURL:  "https://osu.ppy.sh/oauth/token",
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		RedirectURL: os.Getenv("BANCHO_OAUTH_REDIRECT_URL"),
		Scopes:      []string{""},
	}

	// Get the code
	code := c.Query("code")

	if code == "" {
		state, err := genSessionToken(32)
		if err != nil {
			return err
		}

		// make new cookie with state value
		cookie := &fiber.Cookie{
			Name:  "oauth_state",
			Value: state,
		}
		c.Cookie(cookie)

		globals.Logger.WithFields(logrus.Fields{
			"state":    state,
			"redirect": oauthConfig.AuthCodeURL(state),
		}).Debug("Generated new state, redirecting")

		return c.Redirect(oauthConfig.AuthCodeURL(state))
	}

	token, err := getOauth(c, &oauthConfig, code)

	globals.Logger.WithFields(logrus.Fields{
		"token": token,
	}).Info("Got token")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get token",
			"error":   err.Error(),
		})
	}

	// Get the user info
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://osu.ppy.sh/api/v2/me/osu")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get user data",
			"error":   err.Error(),
		})
	}

	defer resp.Body.Close()

	banchoResp := structs.BanchoMe{}
	json.NewDecoder(resp.Body).Decode(&banchoResp)

	user := structs.User{}
	user.BanchoId = banchoResp.Id

	// Check for session token
	sessionToken, _ := globals.CheckAuth(c)
	search := structs.User{
		Session: &structs.Session{
			SessionToken: sessionToken,
		},
	}
	res := structs.User{}

	globals.DBConn.Preload("Session").First(&res, search)

	// Associate ripple accounts if a previous session existed for a bancho user
	if res.RippleId != 0 {
		globals.Logger.WithFields(logrus.Fields{
			"res": res,
		}).Debug("Searched for session and found user")
		user = res
		user.BanchoId = banchoResp.Id
		user.Username = banchoResp.Username

		user.Session.AccessToken = token.AccessToken
		user.Session.RefreshToken = token.RefreshToken

		globals.Logger.WithFields(logrus.Fields{
			"user": user,
		}).Debug("Saving user")

		err = globals.DBConn.Save(&user).Error

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Could not save user",
				"error":   err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:  "session_token",
			Value: user.Session.SessionToken,
		})

		// Clean up after ourselfs
		c.ClearCookie("oauth_state")

		return c.Status(fiber.StatusOK).Redirect("/")

	} else {
		sessionToken, err = genSessionToken(32)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Could not generate session token",
				"error":   err.Error(),
			})
		}
	}

	// check if user with this BanchoId exists
	globals.DBConn.Preload("Session").First(&res, structs.User{BanchoId: user.BanchoId})

	// check if user had a session
	if res.Session.ID != 0 {
		// update the session
		user.Session.SessionToken = sessionToken
		user.Session.AccessToken = token.AccessToken
		user.Session.RefreshToken = token.RefreshToken
		err = globals.DBConn.Save(&user.Session).Error

	} else {
		user.Session = &structs.Session{
			AccessToken:  token.AccessToken,
			SessionToken: sessionToken,
		}
		err = globals.DBConn.Save(&user).Error
	}

	globals.Logger.WithFields(logrus.Fields{
		"user": user,
	}).Debug("Saved User")

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
		Name:  "session_token",
		Value: user.Session.SessionToken,
	})

	// Clean up after ourselfs
	c.ClearCookie("oauth_state")

	return c.Status(fiber.StatusOK).Redirect("/")
}

func getOauth(c *fiber.Ctx, oauthConfig *oauth2.Config, code string) (*oauth2.Token, error) {
	globals.Logger.WithFields(logrus.Fields{
		"code": code,
	}).Debug("Received code")

	// Read oauthState from Cookie
	oauth_state := c.Cookies("oauth_state")

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
