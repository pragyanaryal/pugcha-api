package use_cases

import (
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveBusiness"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/serveUtils"
	"net/http"
)

var (
	BusinessFilterableKeys = []string{"category_id", "user_id", "approved", "approved_by", "blocked", "distance", "near", "type"}
	BusinessSortableKeys   = []string{"created_on", "distance"}
)

func ListBusinesses(temp middleware.AuthUserTemp, r *http.Request) (*[]interface{}, error) {
	params, err := serveUtils.ParseQuery(r)
	if err != nil {
		return nil, err
	}

	params = serveUtils.RefineSortFilter(BusinessFilterableKeys, BusinessSortableKeys, params)

	sql := serveUtils.ToSQL(params, "business")

	service := serveBusiness.BusinessService(repositories.BusinessRepo, repositories.BusinessProfileRepo)

	business, err := service.ListBusinesses(sql)
	if err != nil {
		return nil, err
	}
	return business, nil
}
