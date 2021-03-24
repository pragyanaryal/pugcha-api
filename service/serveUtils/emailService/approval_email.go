package emailService

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/matcornic/hermes/v2"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"golang.org/x/crypto/sha3"
	"strconv"
	"time"
)

// CreateHtmlMessage for approval mail sent by manager
func CreateHtmlMessageApproval(user *map[string]*models.User) (*map[string]string, error) {
	var emails = make(map[string]string, 15)

	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      "Pugcha",
			Link:      "https://example-hermes.com/",
			Logo:      "https://i.ibb.co/zH5z9h6/Logo-p.png",
			Copyright: "Copyright © " + strconv.Itoa(time.Now().Year()) + " Pugcha",
		},
	}

	for key, val := range *user {

		id, token := createToken(val)
		email := hermes.Email{
			Body: hermes.Body{
				Name:   key,
				Intros: []string{"Welcome to Pugcha! We're very excited to have you on board."},
				Actions: []hermes.Action{
					{
						Instructions: "Click the button below to accept the invitation:",
						Button: hermes.Button{
							Color: "#22BC66",
							Text:  "Accept the invitation",
							Link:  config.Configuration.FrontendURL + id + "/" + token,
						},
					},
				},
				Outros: []string{
					"If you did not expect, no further action is required on your part.",
				},
				Signature: "Thanks",
			},
		}
		emailBody, err := h.GenerateHTML(email)
		if err != nil {
			fmt.Println(err, "error here")
			return nil, errors.New("no email")
		}

		emails[val.Email] = emailBody
	}
	return &emails, nil
}

func createToken(user *models.User) (string, string) {
	userId := user.ID
	userPassword := user.Password
	userCreated := user.Status

	secret := userId.String() + string(userPassword) + userCreated
	hash := sha3.Sum224([]byte(secret))

	pass := hex.EncodeToString(hash[:])
	pass = fmt.Sprintf("%x", hash)

	return userId.String(), pass
}
