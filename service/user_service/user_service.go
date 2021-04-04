package user_service

import (
	"github.com/danhper/structomap"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"time"
)

func (users *userService) CreateUser(email string, userType string) (*models.User, error) {
	user := models.User{
		ID:        uuid.New(),
		Email:     email,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Status:    "active",
		Type:      userType,
		Password:  "",
	}

	profileMap := structomap.New().UseSnakeCase().PickAll().Transform(user)

	err := users.userRepo.CreateUser(&profileMap)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (users *userService) CreateAdminUser(email string, userType string) (*models.User, error) {

	pass, err := generatePassword()
	if err != nil {
		return nil, err
	}

	hashPassword, err := HashThePassword(pass)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:        uuid.New(),
		Email:     email,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Status:    "verification_needed",
		Type:      userType,
		Password:  string(hashPassword),
	}
	profileMap := structomap.New().UseSnakeCase().PickAll().Transform(user)

	err = users.userRepo.CreateUser(&profileMap)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (users *userService) ListUsers(filter *serveUtils.SQLFilter) (*[]*json_serializer.UserResponse, error) {
	profile, err := users.userRepo.ListUser(filter)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (users *userService) GetUserByEmail(email string) (*json_serializer.UserResponse, error) {
	user, err := users.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (users *userService) GetUserById(id uuid.UUID) (*json_serializer.UserResponse, error) {
	user, err := users.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (users *userService) UpdateUser(user *models.User) (*models.User, error) {
	err := users.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (users *userService) PatchUser(user uuid.UUID, patch *map[string]interface{}) error {
	err := users.userRepo.PatchUser(user, patch)
	if err != nil {
		return err
	}

	return nil
}

func (users *userService) DeleteUser(user uuid.UUID) error {
	err := users.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}
