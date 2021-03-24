package serveUsers

import (
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(serializer *json_serializer.LoginSerializer) (*json_serializer.UserResponse, error) {
	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)

	service := UserService(userRepo, profileRepo, googleRepo, fbRepo)

	user, err := service.GetUserByEmail(serializer.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(serializer.Password))
	if err != nil {
		return nil, err
	}

	return user, err
}
