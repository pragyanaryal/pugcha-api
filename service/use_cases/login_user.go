package use_cases

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"net/http"
)

func LoginUserOrchestrator(rw http.ResponseWriter, userID uuid.UUID, level string, from string,
	serializer *json_serializer.LoginSerializer) (error, *json_serializer.UserResponse, *json_serializer.TokenSerializer) {

	if from == "google" || from == "facebook" {
		token, err := serveUtils.CreateToken(userID, level, from, 0)
		if err != nil {
			return err, nil, nil
		}

		return nil, nil, token
	}

	err := validateLoginRequest(serializer)
	if err != nil {
		return err, nil, nil
	}

	user, err := serveUsers.Authenticate(serializer)
	if err != nil {
		return err, nil, nil
	}

	token, err := serveUtils.CreateToken(user.ID, user.Type, "normal", 0)
	//serveUtils.SetCookie(rw, token)

	return nil, user, token
}

func validateLoginRequest(a interface{}) error {
	validate := validator.New()
	return validate.Struct(a)
}
