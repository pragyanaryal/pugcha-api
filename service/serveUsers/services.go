package serveUsers

import (
	"github.com/sethvargo/go-password/password"
	"gitlab.com/ProtectIdentity/pugcha-backend/models/repositories_interface"
	"golang.org/x/crypto/bcrypt"
)

var (
	GoogleService = googleService{}
)

type googleService struct{}
type facebookService struct{}

type (
	profileRepo repositories_interface.UserProfileRepository
	userRepo    repositories_interface.UserRepository
	fbRepo      repositories_interface.FacebookProfile
	googleRepo  repositories_interface.GoogleProfile
)

type userService struct {
	userRepo        userRepo
	userProfile     profileRepo
	googleProfile   googleRepo
	facebookProfile fbRepo
}

func UserService(userRepo userRepo, profile profileRepo, googleProfile googleRepo, fbProfile fbRepo) *userService {
	return &userService{
		userRepo:        userRepo,
		userProfile:     profile,
		googleProfile:   googleProfile,
		facebookProfile: fbProfile,
	}
}

func generatePassword() (string, error) {
	pass, err := password.Generate(24, 10, 10, false, false)
	if err != nil {
		return "", err
	}

	return pass, nil
}

func HashThePassword(plainPassword string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
