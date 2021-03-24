package use_cases

import (
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUsers"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"net/http"
)

var (
	UserFilterableKeys = []string{"email", "status", "type", "gender", "contact", "first_name", "last_name"}
	UserSortableKeys   = []string{"created_on", "updated_on", "gender", "type"}
)

func ListUsers(temp middleware.AuthUserTemp, r *http.Request) (*[]*json_serializer.UserResponse, error) {
	params, err := serveUtils.ParseQuery(r)
	if err != nil {
		return nil, err
	}

	params = serveUtils.RefineSortFilter(UserFilterableKeys, UserSortableKeys, params)

	sql := serveUtils.ToSQL(params, "user")

	var (
		userRepo    = repositories.UserRepo
		profileRepo = repositories.UserProfileRepo
		googleRepo  = repositories.GoogleRepo
		fbRepo      = repositories.FacebookRepo
		service     = serveUsers.UserService(userRepo, profileRepo, googleRepo, fbRepo)
	)

	user, err := service.ListUsers(sql)
	if err != nil {
		return nil, err
	}
	return user, nil
}
