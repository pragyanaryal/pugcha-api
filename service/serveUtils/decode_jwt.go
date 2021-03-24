package serveUtils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
)

func DecodeJWT(tokens string, key string) (*jwt.MapClaims, error) {

	token, err := jwt.Parse(tokens, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if key == "access" {
			return []byte(config.Configuration.AccessKey), nil
		} else {
			return []byte(config.Configuration.RefreshKey), nil
		}
	})

	if err != nil {
		return nil, err
	} else {
		claims, ok := token.Claims.(jwt.MapClaims)

		if ok && token.Valid {
			return &claims, nil
		}
		return nil, err
	}
}
