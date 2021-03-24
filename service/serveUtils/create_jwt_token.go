package serveUtils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"time"
)

//const AccessKey = "KoTYY2DDOLjh/IvApUH4cR+wt5rLdcGveEBrwSo75CJg+579rf2whgh/M3N6Pugqy+C6cA0Loh2o8RI5cpw89A=="
//const RefreshKey = "XX7P/+XaPNX5JklESj1LihxUIat7Yq/L50Cz+VOLcewUSpGNGoqyWwkwbUCRXAMFyCAF+yL1hRPMO23juguXOg=="

func CreateToken(userId uuid.UUID, level string, from string, times int) (*json_serializer.TokenSerializer, error) {
	accessToken, acID, err := CreateJWTAccessToken(userId, level, from)
	if err != nil {
		return nil, err
	}

	refreshToken, err := CreateJWTRefreshToken(userId, level, from, acID, times)
	if err != nil {
		return nil, err
	}

	token := json_serializer.TokenSerializer{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &token, nil
}

func CreateJWTAccessToken(userId uuid.UUID, level string, from string) (string, uuid.UUID, error) {
	atClaims := jwt.MapClaims{}

	acId := uuid.New()

	atClaims["authorized"] = true
	atClaims["access_id"] = acId
	atClaims["userId"] = userId
	atClaims["type"] = level
	atClaims["loggedFrom"] = from
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	token, err := at.SignedString([]byte(config.Configuration.AccessKey))
	if err != nil {
		return "", [16]byte{}, err
	}

	return token, acId, nil
}

func CreateJWTRefreshToken(usId uuid.UUID, lev string, fr string, acId uuid.UUID, t int) (string, error) {
	atClaims := jwt.MapClaims{}

	atClaims["refresh_id"] = uuid.New()
	atClaims["access_id"] = acId
	atClaims["userId"] = usId
	atClaims["times"] = t
	atClaims["type"] = lev
	atClaims["loggedFrom"] = fr
	if lev == "user" {
		atClaims["exp"] = time.Now().Add(time.Hour * 24 * 3).Unix()
	} else {
		atClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	token, err := at.SignedString([]byte(config.Configuration.RefreshKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
