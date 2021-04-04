package business_service

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/form"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"gitlab.com/ProtectIdentity/pugcha-backend/models/repositories_interface"
	"gitlab.com/ProtectIdentity/pugcha-backend/serializer/json_serializer"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type businessService struct {
	businessRepo repositories_interface.BusinessRepository
	profileRepo  repositories_interface.BusinessProfileRepository
}

func BusinessService(
	businessRepo repositories_interface.BusinessRepository,
	profile repositories_interface.BusinessProfileRepository) *businessService {
	return &businessService{businessRepo: businessRepo, profileRepo: profile}
}

type OpeningHours struct {
	Day         string `json:"day,omitempty"`
	OpeningTime string `json:"opening_time,omitempty"`
	ClosingTime string `json:"closing_time,omitempty"`
	Opened      bool   `json:"opened,omitempty"`
}

type Location struct {
	Latitude  string
	Longitude string
}

var DecodeUUID = func(values []string) (interface{}, error) {
	return uuid.Parse(values[0])
}

var DecodeLocation = func(values []string) (interface{}, error) {
	location := Location{}
	_ = json.Unmarshal([]byte(values[0]), &location)
	if location.Latitude == "" || location.Longitude == "" {
		return models.Location{Latitude: nil, Longitude: nil},
			errors.New("empty")
	}
	lat, err := strconv.ParseFloat(location.Latitude, 64)
	lon, err1 := strconv.ParseFloat(location.Longitude, 64)

	if err != nil || err1 != nil || lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return models.Location{Latitude: nil, Longitude: nil},
			errors.New("invalid value")
	}
	return models.Location{
		Latitude:  &lat,
		Longitude: &lon,
	}, nil
}

func contains(check string) (int, error) {
	Days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	for key, val := range Days {
		if check == val {
			return key + 1, nil
		}
	}
	return 0, errors.New("not found")
}

var DecodeOpeningHours = func(values []string) (interface{}, error) {
	var openingHour []models.OpeningHours
	for _, val := range values {
		opening := OpeningHours{}
		err := json.Unmarshal([]byte(val), &opening)
		if err != nil {
			return nil, err
		}

		loca, err := time.LoadLocation("Asia/Kathmandu")

		open, err := time.ParseInLocation("15:04", opening.OpeningTime, loca)
		closes, err1 := time.ParseInLocation("15:04", opening.ClosingTime, loca)

		if err != nil || err1 != nil {
			return []models.OpeningHours{}, errors.New("invalid value")
		}

		opened := open.Format("15:04 -0700")
		closed := closes.Format("15:04 -0700")

		var day models.Weekday

		temp, err := contains(strings.Title(strings.ToLower(opening.Day)))
		if temp >= 1 && temp <= 8 {
			day = models.Weekday(temp)
		}

		if day < 1 || day > 7 {
			return []models.OpeningHours{}, errors.New("invalid day")
		}

		openingHour = append(openingHour, models.OpeningHours{
			WeekDay:     day,
			OpeningTime: &opened,
			ClosingTime: &closed,
			Opened:      true,
		})
	}
	return openingHour, nil
}

var DecodeNullTime = func(values []string) (interface{}, error) {
	if len(values[0]) == 4{
		values[0] = values[0] + "-01-01"
	}
	establishedDate, err := time.Parse("2006-01-02", values[0])
	if err != nil {
		return nil, err
	}
	return &establishedDate, nil
}

var DecodeAddresses = func(values []string) (interface{}, error) {
	var address []models.Address
	for _, val := range values {
		add := models.Address{}
		err := json.Unmarshal([]byte(val), &add)
		if err != nil {
			return nil, err
		}

		address = append(address, add)
	}
	return address, nil
}

var DecodeOwner = func(values []string) (interface{}, error) {
	var owner []models.Owner
	for _, val := range values {
		own := models.Owner{}
		err := json.Unmarshal([]byte(val), &own)
		if err != nil {
			return nil, err
		}
		owner = append(owner, own)
	}
	return owner, nil
}

func GetBusinessProfileData(r *http.Request, uid uuid.UUID, image []string) (*models.BusinessProfile, error) {

	var business models.BusinessProfile

	var decoder *form.Decoder
	decoder = form.NewDecoder()

	decoder.SetMaxArraySize(15)

	// Decodes the uuid Type of variable sent by the user
	decoder.RegisterCustomTypeFunc(DecodeUUID, uuid.UUID{})

	// Decodes the Location which involves latitude and longitude, which is sent by the user as a json object
	decoder.RegisterCustomTypeFunc(DecodeLocation, models.Location{})

	// Decodes the opening hour time format, user send this as an array of json object. The opening hour may consist of null time
	decoder.RegisterCustomTypeFunc(DecodeOpeningHours, []models.OpeningHours{})

	// Decodes the Null time type, where the user sends the date type which can be null
	decoder.RegisterCustomTypeFunc(DecodeNullTime, &time.Time{})

	// Decodes the address type , where user sends it as a array of json object
	decoder.RegisterCustomTypeFunc(DecodeAddresses, []models.Address{})

	// Decodes the owner type, if the user has sent any. The user sends it as an array of json object.
	decoder.RegisterCustomTypeFunc(DecodeOwner, []models.Owner{})

	err := decoder.Decode(&business, r.PostForm)
	if err != nil {
		return nil, err
	}

	business.UserId = uid
	business.Picture = image

	return &business, nil
}

func PopulateBusinessId(id uuid.UUID, profile *models.BusinessProfile) (*models.BusinessProfile, error) {
	// Update owner
	var owner []models.Owner
	for _, val := range profile.Owner {
		val.BusinessId = id
		val.Id = uuid.New()
		owner = append(owner, val)
	}

	// Update Address
	var address []models.Address
	for _, val := range profile.Address {
		val.BusinessId = id
		val.Id = uuid.New()
		address = append(address, val)
	}

	// Update Opening Times
	var opening []models.OpeningHours
	for _, val := range profile.OpeningHours {
		val.BusinessId = id
		opening = append(opening, val)
	}

	profile.Owner = owner
	profile.OpeningHours = opening
	profile.Address = address
	profile.BusinessId = id

	return profile, nil
}

func PopulateBusinessIdOnUpdate(id uuid.UUID, profile *json_serializer.UpdateAbleBusiness) (*json_serializer.UpdateAbleBusiness, error) {
	// Update owner
	if profile.Owner != nil {
		var owner []models.Owner
		for _, val := range *profile.Owner {
			val.BusinessId = id
			val.Id = uuid.New()
			owner = append(owner, val)
		}
		profile.Owner = &owner
	}

	// Update Address

	if profile.Address != nil {
		var address []models.Address
		for _, val := range *profile.Address {
			val.BusinessId = id
			val.Id = uuid.New()
			address = append(address, val)
		}
		profile.Address = &address
	}

	// Update Opening Times
	if profile.OpeningHours != nil {
		var opening []models.OpeningHours
		for _, val := range *profile.OpeningHours {
			val.BusinessId = id
			opening = append(opening, val)
		}
		profile.OpeningHours = &opening
	}

	return profile, nil
}
