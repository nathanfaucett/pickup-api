package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aicacia/pickup/app/config"
	"github.com/aicacia/pickup/app/jwt"
	"github.com/aicacia/pickup/app/model"
	"github.com/aicacia/pickup/app/repository"
	"github.com/aicacia/pickup/app/util"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  string
)

// GetProviderRedirect
//
//	@Summary		Redirect to provider
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string	true	"Provider"
//	@Success		302
//	@Router			/oauth2/{provider} [get]
func GetProviderRedirect(c *fiber.Ctx) error {
	var url string
	switch c.Params("provider") {
	case "google":
		url = googleOauthConfig.AuthCodeURL(oauthStateString)
	default:
		return model.NewError(400).Send(c)
	}
	return c.Redirect(url, http.StatusFound)
}

// GetProviderCallback
//
//	@Summary		Redirects with provider's token
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string	true	"provider"
//	@Param			state	    path		string	true	"state"
//	@Param			code	    path		string	true	"code"
//	@Success		302
//	@Router			/oauth2/{provider}/callback [get]
func GetProviderCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	userInfo, err := getUserInfo(provider, c.FormValue("state"), c.FormValue("code"))
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		errorValue := model.NewError(500).AddError("application", "internal", provider)
		return RedirectWithError(c, errorValue)
	}
	if !userInfo.IsEmailVerified() {
		errorValue := model.NewError(400).AddError("email", "unverified", provider)
		return RedirectWithError(c, errorValue)
	}
	user, err := repository.GetUserByEmail(userInfo.GetEmail())
	if err != nil {
		log.Printf("Error getting user: %v", err)
		errorValue := model.NewError(500).AddError("application", "internal", "database")
		return RedirectWithError(c, errorValue)
	}
	if user == nil {
		newUser, err := repository.CreateUserFromEmail(userInfo.GetEmail())
		if err != nil {
			log.Printf("Error creating user: %v", err)
			errorValue := model.NewError(500).AddError("application", "internal", "database")
			return RedirectWithError(c, errorValue)
		}
		user = &newUser
	}
	now := time.Now().UTC()
	claims := jwt.Claims{
		Subject:          user.Id,
		NotBeforeSeconds: now.Unix(),
		IssuedAtSeconds:  now.Unix(),
		ExpiresAtSeconds: now.Unix() + int64(config.Get().JWT.Expires.Seconds),
	}
	token, err := jwt.CreateToken(&claims)
	if err != nil {
		log.Printf("Erro creating token: %v", err)
		errorValue := model.NewError(500).AddError("application", "internal", "token")
		return RedirectWithError(c, errorValue)
	}
	return c.Redirect(fmt.Sprintf("%s/oauth/callback?token=%s", config.Get().UI.URI, url.QueryEscape(token)), http.StatusFound)
}

func RedirectWithError(c *fiber.Ctx, errorValue *model.ErrorST) error {
	bytes, err := json.Marshal(errorValue)
	if err != nil {
		log.Printf("Error converting error to json: %v", err)
		bytes = []byte("{\"errors\":[]}")
	}
	return c.Redirect(fmt.Sprintf("%s/oauth/callback?error=%s", config.Get().UI.URI, url.QueryEscape(string(bytes))), http.StatusFound)
}

type UserInfo interface {
	GetEmail() string
	IsEmailVerified() bool
}

func getUserInfo(provider, state, code string) (UserInfo, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	switch provider {
	case "google":
		return getGoogleUserInfo(code)
	default:
		return nil, fmt.Errorf("invalid provider")
	}
}

type GoogleUserInfo struct {
	Id            string `json:"id" validate:"required"`
	Email         string `json:"email" validate:"required"`
	VerifiedEmail bool   `json:"verified_email" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Picture       string `json:"picture" validate:"required"`
}

func (user *GoogleUserInfo) GetEmail() string {
	return user.Email
}

func (user *GoogleUserInfo) IsEmailVerified() bool {
	return user.VerifiedEmail
}

func getGoogleUserInfo(code string) (*GoogleUserInfo, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var user GoogleUserInfo
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func InitOAuth2() (err error) {
	oauthStateString, err = util.GenerateRandomHex(32)
	if err != nil {
		return
	}
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.Get().URI + "/oauth2/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return
}
