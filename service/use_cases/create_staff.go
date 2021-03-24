package use_cases

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils/emailService"
)

func CreateStaffUserOrchestrator(staff *json_serializer.CreateAdminUserRequest) (*json_serializer.UserResponse, error) {
	err := validateStaffCreationRequest(staff)
	if err != nil {
		return nil, err
	}

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)

	service := serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)

	old, err := service.GetUserByEmail(staff.Email)
	if err != nil {
		user, err := service.CreateAdminUser(staff.Email, "admin")
		if err != nil {
			return nil, err
		}

		users := make(map[string]*models.User)
		users[staff.FirstName+" "+staff.LastName] = user

		err = emailService.SendEmail("approval", &users, false)

		if err != nil {
			return nil, err
		}

		new, _ := service.GetUserByEmail(staff.Email)
		return new, nil
	}
	return old, errors.New("an user with email already exists")
}

func validateStaffCreationRequest(a interface{}) error {
	validate := validator.New()
	err := validate.StructPartial(a)
	if err != nil {
		return errors.New("validation error")
	}
	return nil
}
