package serveUsers

import "golang.org/x/oauth2"

var facebookOauthConfig = &oauth2.Config{
	ClientID:     "",
	ClientSecret: "",
	Endpoint:     oauth2.Endpoint{},
	RedirectURL:  "",
	Scopes:       nil,
}
