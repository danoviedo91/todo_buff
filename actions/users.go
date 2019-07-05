package actions

import (
	"log"
	"strings"

	"github.com/danoviedo91/todo_buff/mailers"
	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
)

// UsersNew ...
func UsersNew(c buffalo.Context) error {

	if uid := c.Session().Get("current_user_id"); uid == nil {
		if err := c.Redirect(301, "/"); err != nil {
			log.Fatal(err)
		}
	}

	u := models.User{}
	c.Set("user", u)
	c.Set("context", c) // Used for HrefCancelBtn helper
	log.Println(c.Session().Get("filterStatus"))
	return c.Render(200, r.HTML("users/new.html"))

}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}

	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	// Generate a password that is 12 characters long with 4 digits, 4 symbols,
	// allowing upper and lower case letters, allowing repeat characters.
	pwd, err := password.Generate(12, 4, 4, false, true)

	if err != nil {
		log.Fatal(err)
	}

	u.Password = pwd
	u.PasswordConfirmation = u.Password

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		log.Println(verrs.Errors)
		c.Set("user", u)
		c.Set("errors", verrs)
		c.Set("context", c) // Used for HrefCancelBtn helper
		return c.Render(200, r.HTML("users/new.html"))
	}

	adminUserID := c.Value("current_user").(*models.User).ID.String()

	emailBody, err := mailers.GenerateEmailHTMLBody(tx, u, adminUserID)

	if err != nil {
		errors.WithStack(err)
	}

	reader := strings.NewReader(emailBody)
	browser.OpenReader(reader)

	// from := mail.NewEmail("Example User", "doviedo@wawand.co")
	// subject := "Sending with Twilio SendGrid is Fun"
	// to := mail.NewEmail("Example User", "doviedo@wawand.co")
	// plainTextContent := "and easy to do anywhere, even with Go"
	// htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	// response, err := client.Send(message)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// fmt.Println(response.StatusCode)
	// fmt.Println(response.Body)
	// fmt.Println(response.Headers)

	return c.Redirect(302, "/")
}

// UsersEdit ...
func UsersEdit(c buffalo.Context) error {
	return nil
}

// UsersUpdate ...
func UsersUpdate(c buffalo.Context) error {
	return nil
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/signin")
		}
		return next(c)
	}
}
