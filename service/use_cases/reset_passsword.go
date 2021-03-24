package use_cases

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"golang.org/x/crypto/sha3"
	"time"
	"unicode"
)

// ResetPasswordOrchestrator ...
func ResetPasswordOrchestrator(userId uuid.UUID, token string,
	password *json_serializer.PasswordSerializer, kind string) (string, *json_serializer.UserResponse, error) {

	err := validateNewPassword(password)
	if err != nil {
		return "", nil, err
	}

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)
	service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)

	user, errs := service.GetUserById(userId)
	if errs != nil {
		return "", nil, errs
	}

	userd := user.ID
	userPassword := user.Password
	userCreated := user.Status

	secret := userd.String() + userPassword + userCreated

	hash := sha3.Sum224([]byte(secret))
	pass := hex.EncodeToString(hash[:])
	pass = fmt.Sprintf("%x", hash)

	if token != pass {
		return "", nil, errors.New("the link is tampered")
	}

	thePassword, er := serveUsers.HashThePassword(password.Password)
	if er != nil {
		return "", nil, err
	}

	if kind == "reset" {
		patch := make(map[string]interface{})
		patch["password"] = thePassword
		patch["updated_on"] = time.Now()

		err = service.PatchUser(user.ID, &patch)
		if err != nil {
			return "", nil, err
		}
		return "", nil, nil
	}
	return string(thePassword), user, nil
}

func validateNewPassword(a *json_serializer.PasswordSerializer) error {
	validate := validator.New()

	err := validate.Struct(a)
	if err != nil {
		return err
	}

	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range a.Password {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		return errors.New("error")
	}
	return nil
}
