package grifts

import (
	"github.com/danoviedo91/todo_buff/models"
	"github.com/gofrs/uuid"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("createadmin", "Creates Admin Account")
	grift.Add("createadmin", func(c *grift.Context) error {

		admin := models.User{}

		hash, err := bcrypt.GenerateFromPassword([]byte("banana"), bcrypt.DefaultCost)

		if err != nil {
			return errors.WithStack(err)
		}

		admin.ID = uuid.Must(uuid.NewV4())
		admin.FirstName = "Daniel"
		admin.LastName = "Oviedo"
		admin.Email = "admin@wawand.co"
		admin.IsAdmin = true
		admin.PasswordHash = string(hash)

		err = models.DB.Create(&admin)

		return err
	})

	grift.Desc("deleteadmin", "Deletes Admin Account")
	grift.Add("deleteadmin", func(c *grift.Context) error {

		admin := models.User{}

		if err := models.DB.Where("email = ?", "admin@wawand.co").First(&admin); err != nil {
			return errors.WithStack(err)
		}

		if err := models.DB.Destroy(&admin); err != nil {
			return errors.WithStack(err)
		}

		return nil
	})

})
