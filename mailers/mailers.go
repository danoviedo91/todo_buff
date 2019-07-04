package mailers

import (
	"fmt"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/pop"
	"github.com/matcornic/hermes"
	"github.com/pkg/errors"
)

// GenerateEmailHTMLBody ...
func GenerateEmailHTMLBody(tx *pop.Connection, u *models.User, adminUserID string) (emailBody string, err error) {

	// Get admin info

	adminUser := models.User{}

	if err := tx.Find(&adminUser, adminUserID); err != nil {
		return "", errors.WithStack(err)
	}

	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: fmt.Sprintf("%v %v", adminUser.FirstName, adminUser.LastName),
			Link: "https://example-hermes.com/",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: fmt.Sprintf("%v %v", u.FirstName, u.LastName),
			Intros: []string{
				"Welcome to ToDo Dashboard App! We're very excited to have you on board. Your login information is:",
			},
			Dictionary: []hermes.Entry{
				{
					Key:   "Email",
					Value: u.Email,
				},
				{
					Key:   "Password",
					Value: u.Password,
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err = h.GenerateHTML(email)

	if err != nil {
		return "", err
	}

	return emailBody, nil

}
