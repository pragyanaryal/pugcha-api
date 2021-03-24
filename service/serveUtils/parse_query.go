package serveUtils

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-playground/form"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
)

var (
	AllBusinessSelect = []string{"picture", "business_id", "contact", "location", "description",
		"ST_Distance(st_transform(geom, 4326), ref_geom) AS distance", "category_name", "business_name",
		"pan_number", "vat_number", "email", "website", "established_date", "address", "opening"}
	SmallBusinessSelect = []string{"picture", "business_id", "contact", "location", "description", "category_name",
		"ST_Distance(st_transform(geom, 4326), ref_geom) AS distance", "business_name"}
	AllBusinessWithoutDistance = []string{"picture", "business_id", "contact", "location", "description",
		"category_name", "business_name", "pan_number", "vat_number", "email", "website",
		"established_date", "address", "opening"}
	SmallBusinessWithoutDistance = []string{"picture", "business_id", "contact", "location", "description", "business_name",
		"category_name"}
)

type SQLFilter struct {
	Limit  uint16
	Offset uint16
	Where  []string
	Sort   map[string]string
	Filter map[string]string
	Args   map[string][]interface{}
	Amount string
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func RefineSortFilter(field []string, sort []string, params *json_serializer.FilterRequest) *json_serializer.FilterRequest {
	for key := range params.Fields {
		present := Contains(field, key)
		if !present {
			delete(params.Fields, key)
		}
	}

	var newSort []string
	for _, val := range params.Sort {
		vales := val[1:]
		present := Contains(sort, vales)
		if present {
			newSort = append(newSort, val)
		}
	}
	params.Sort = newSort
	return params
}

func isJSON(s string) (map[string]interface{}, bool) {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(s), &js)
	if err != nil {
		return nil, false
	}
	return js, true
}

func GetFieldsNew(strings []string) (interface{}, error) {
	var temp map[string]interface{}
	err := json.Unmarshal([]byte(strings[0]), &temp)
	if err != nil {
		return nil, err
	}

	p, err := DecodeFields(temp)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func ParseQuery(r *http.Request) (*json_serializer.FilterRequest, error) {
	var filter json_serializer.FilterRequest

	var decoder *form.Decoder
	decoder = form.NewDecoder()
	decoder.RegisterCustomTypeFunc(GetFieldsNew, map[string]*json_serializer.Field{})

	err := decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		return nil, err
	}

	return &filter, nil
}

func DecodeFields(temp map[string]interface{}) (map[string]*json_serializer.Field, error) {
	fields := map[string]*json_serializer.Field{}

	for key, val := range temp {
		f, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		m, check := isJSON(string(f))

		if check {
			var t []interface{}
			for _, vald := range m {
				t = append(t, vald)
			}

			if key == "near" {
				for keys, vals := range m {
					if keys == "latitude" {
						t[0] = vals
					} else {
						t[1] = vals
					}
				}
			}

			temp := json_serializer.Field{
				Name:       key,
				FilterOp:   "eq",
				FilterData: t,
			}
			fields[key] = &temp
		} else {
			var t []interface{}
			temp := json_serializer.Field{
				Name:       key,
				FilterOp:   "eq",
				FilterData: append(t, val),
			}
			fields[key] = &temp
		}
	}
	return fields, nil
}

func ToSQL(params *json_serializer.FilterRequest, kind string) *SQLFilter {
	sort := map[string]string{}

	var (
		filter map[string]string
		args   map[string][]interface{}
		limit  uint16
		offset uint16
		where  []string
		amount string
	)

	for _, val := range params.Sort {
		sign := val[0:1]
		if sign == "+" {
			temp := val[1:] + " ASC"
			sort[val[1:]] = temp
		} else if sign == "-" {
			temp := val[1:] + " DESC"
			sort[val[1:]] = temp
		} else {
			temp := val + " ASC"
			sort[val] = temp
		}
	}

	if kind == "business" {
		filter, args, where, amount = BusinessFilterToSql(params.Fields)
	} else {
		filter, args, where, amount = UserFilterToSql(params.Fields)
	}

	limit = params.Limit
	offset = params.Offset

	if params.Limit <= 0 {
		limit = 10
	}
	if params.Offset < 0 {
		offset = 0
	}

	sql := SQLFilter{
		Limit:  limit,
		Offset: offset,
		Where:  where,
		Sort:   sort,
		Filter: filter,
		Args:   args,
		Amount: amount,
	}

	return &sql
}

func UserFilterToSql(params map[string]*json_serializer.Field) (map[string]string, map[string][]interface{}, []string, string) {
	filter := map[string]string{}
	args := map[string][]interface{}{}
	var where []string

	for key, val := range params {
		if key == "email" {
			var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if emailRegex.MatchString(val.FilterData[0].(string)) {
				filter[key] = key + " = ?"
				args[key] = val.FilterData
			}
		} else if key == "status" {
			temp := val.FilterData[0].(string)
			if temp == "active" || temp == "verification_needed" || temp == "approval_needed" || temp == "blocked" {
				filter[key] = key + " = ?"
				args[key] = val.FilterData
			}
		} else if key == "type" {
			temp := val.FilterData[0].(string)
			if temp == "user" || temp == "admin" || temp == "super" {
				filter[key] = key + " = ?"
				args[key] = val.FilterData
			}
		} else if key == "gender" {
			temp := val.FilterData[0].(string)
			if temp == "male" || temp == "female" || temp == "rest" {
				filter[key] = key + " = ?"
				args[key] = val.FilterData
			}
		} else {
			filter[key] = key + " = ?"
			args[key] = val.FilterData
		}
	}

	return filter, args, where, ""
}

func BusinessFilterToSql(params map[string]*json_serializer.Field) (map[string]string, map[string][]interface{}, []string, string) {

	filter := map[string]string{}
	args := map[string][]interface{}{}
	var where []string
	var amount string

	if val, ok := params["distance"]; ok {
		filter["distance"] = "distance" + " < ? "
		args["distance"] = val.FilterData
	}

	if _, ok := params["near"]; ok {
		if val, ok1 := params["type"]; ok1 {
			if val.FilterData[0].(string) == "all" {
				where = AllBusinessSelect
				amount = "all"
			} else {
				where = SmallBusinessSelect
				amount = "small"
			}
		} else {
			where = SmallBusinessSelect
			amount = "small"
		}
	} else {
		if val, ok1 := params["type"]; ok1 {
			if val.FilterData[0].(string) == "all" {
				where = AllBusinessWithoutDistance
				amount = "all"
			} else {
				where = SmallBusinessWithoutDistance
				amount = "small"
			}
		} else {
			where = SmallBusinessWithoutDistance
			amount = "small"
		}
	}

	for key, val := range params {
		if key != "distance" && key != "type" {
			if key == "category_id" || key == "approved_by" || key == "user_id" {
				_, err := uuid.Parse(val.FilterData[0].(string))
				if err == nil {
					filter[key] = key + " = ?"
					args[key] = val.FilterData
				}
			} else if key == "approved" || key == "blocked" {
				if val.FilterData[0].(bool) == false || val.FilterData[0].(bool) == true {
					filter[key] = key + " = ?"
					args[key] = val.FilterData
				}
			} else {
				filter[key] = key + " = ?"
				args[key] = val.FilterData
			}
		}
	}

	return filter, args, where, amount
}
