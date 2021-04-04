package business

import (
	"net/http"

	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/business_service"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
)

var (
	FilterableKeys = []string{"category_id", "user_id", "approved", "approved_by", "blocked", "distance", "near", "type"}
	SortableKeys   = []string{"created_on", "distance"}
)

func ListBusinesses(temp middleware.AuthUserTemp, r *http.Request) (*[]interface{}, error) {
	params, err := serveUtils.ParseQuery(r)
	if err != nil {
		return nil, err
	}

	params = serveUtils.RefineSortFilter(FilterableKeys, SortableKeys, params)

	sql := serveUtils.ToSQL(params, "business")

	service := business_service.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	business, err := service.ListBusinesses(sql)
	if err != nil {
		return nil, err
	}
	return business, nil
}
