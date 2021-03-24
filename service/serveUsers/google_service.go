package serveUsers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/log"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"time"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     config.Configuration.GoogleClient,
	ClientSecret: config.Configuration.GoogleSecret,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/user.birthday.read",
	},
	Endpoint: google.Endpoint,
}

func (users *googleService) GoogleLogin(rw http.ResponseWriter, r *http.Request) {

	oauthState := users.generateStateOauthCookie(rw)
	u := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)

	http.Redirect(rw, r, u, http.StatusTemporaryRedirect)
}

func (users *googleService) generateStateOauthCookie(rw http.ResponseWriter) string {

	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}

	http.SetCookie(rw, &cookie)

	return state
}

// ................................................................................. //
// CallBack .................................................................  CallBack
// ................................................................................. //

func (users *googleService) GoogleCallback(r *http.Request) (*oauth2.Token, *json_serializer.GoogleResponse, error) {

	// Reading oauthState from Cookie
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, nil, err
	}

	if r.FormValue("state") != oauthState.Value {
		log.Info().Msg("invalid oauth google state")
		return nil, nil, err
	}

	token, response, err := users.getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, nil, err
	}

	var data json_serializer.GoogleResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, nil, err
	}

	return token, &data, nil
}

func (users *googleService) getUserDataFromGoogle(code string) (*oauth2.Token, []byte, error) {

	// Using code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	return token, contents, nil
}
