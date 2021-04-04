package user

import (
	"time"
	
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/user_service"
)

func VerifyAccountOrchestrator(userId uuid.UUID, token string, password *json_serializer.PasswordSerializer) (
	*json_serializer.UserResponse, error) {

	pass, user, err := ResetPasswordOrchestrator(userId, token, password, "verify")
	if err != nil {
		return nil, err
	}

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
	)
	service := user_service.UserService(userRepo, profileRepo, googleRepo, fbRepo)

	patch := make(map[string]interface{})
	patch["password"] = pass
	patch["updated_on"] = time.Now()
	patch["status"] = "approval_needed"

	err = service.PatchUser(user.ID, &patch)
	if err != nil {
		return nil, err
	}

	new, _ := service.GetUserById(userId)
	return new, nil
}
