package routes

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxkruse/magnusopus/backend/globals"
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/maxkruse/magnusopus/backend/utils"
	"golang.org/x/oauth2"
)

func GetOAuthRipple(c *fiber.Ctx) error {
	sess, err := globals.SessionStore.Get(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
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
		sess.Regenerate()
		id := sess.ID()
		sess.Set("oauth_state", id)
		// Save session
		if err := sess.Save(); err != nil {
			return err
		}

		return c.Redirect(oauthConfig.AuthCodeURL(id))
	}

	token, err := getOauth(c, &oauthConfig, code)

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	rippleResp := structs.RippleSelf{}

	client := oauthConfig.Client(context.Background(), token)

	// Get the user info
	resp, err := client.Get("https://ripple.moe/api/v1/users/self")

	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}
	json.NewDecoder(resp.Body).Decode(&rippleResp)

	// Save user
	var user structs.User

	localDB := globals.DBConn
	localDB.Preload("Session").First(&user, "ripple_id = ?", rippleResp.UserId)

	user.RippleId = rippleResp.UserId
	user.Username = rippleResp.Username
	user.Sessions = append(user.Sessions, structs.Session{SessionToken: sess.ID()})

	localDB.Save(&user)

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
	// Read oauthState from Cookie
	oauth_state := sess.Get("oauth_state")

	if c.Query("state") != oauth_state {
		return nil, fiber.ErrBadRequest
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return token, nil
}

func Logout(c *fiber.Ctx) error {
	sess, err := globals.SessionStore.Get(c)
	if err != nil {
		return err
	}

	self, err := utils.GetSelf(c)
	if err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	// Delete session for self
	localDB := globals.DBConn
	localDB.Delete(self.Sessions, "user_id = ? ", self.ID)

	// Logout the user by invalidating the cookie
	sess.Destroy()

	if err := sess.Save(); err != nil {
		return utils.DefaultErrorMessage(c, err, fiber.StatusInternalServerError)
	}

	return c.Redirect("/")
}
