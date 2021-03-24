package serveUtils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"net/http"
)

func RefreshToken(r *http.Request) (*json_serializer.TokenSerializer, error) {
	access_cookie, err := r.Cookie("pugcha_token")
	if err != nil {
		return nil, err
	}

	accClaims, accErrs := parseUnverified(access_cookie.Value)
	if accErrs != nil {
		return nil, err
	}

	refresh_cookie, err := r.Cookie("pugcha_refresh")
	if err != nil {
		return nil, err
	}

	refClaims, refErrs := DecodeJWT(refresh_cookie.Value, "refresh")
	if refErrs != nil {
		return nil, err
	}

	accId := (*accClaims)["access_id"].(string)
	accIdRefresh := (*refClaims)["access_id"].(string)

	if accId == accIdRefresh {
		times := int((*refClaims)["times"].(float64))

		if times < 80 {
			usId := uuid.MustParse((*accClaims)["userId"].(string))
			level := (*accClaims)["type"].(string)
			from := (*accClaims)["loggedFrom"].(string)

			token, err := CreateToken(usId, level, from, times+1)
			if err != nil {
				fmt.Println("4")
				return nil, err
			}
			return token, nil
		}
	}
	return nil, err
}

func parseUnverified(tokens string) (*jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokens, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return &claims, nil
	}
	return nil, err
}
