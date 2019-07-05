package grifts

import (
	"math/rand"
	"time"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gofrs/uuid"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	"github.com/wawandco/fako"
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

	grift.Desc("createtodoes", "Creates 1000 random todoes for admin@wawand.co")
	grift.Add("createtodoes", func(c *grift.Context) error {

		numberOfTodoes := 1000

		for i := 0; i < numberOfTodoes; i++ {

			todo := models.Todo{}

			// Generates random Title and Description
			fako.Fill(&todo)

			// Generate random timestamp (time.Time)
			randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
			randomNow := time.Unix(randomTime, 0)

			todo.Deadline = randomNow

			// Get user.ID

			user := models.User{}

			if err := models.DB.Where("email = ?", "admin@wawand.co").First(&user); err != nil {
				return errors.WithStack(err)
			}

			todo.UserID = user.ID

			// All tasks are NOT completed...

			todo.Completed = false

			// Create todo

			if err := models.DB.Create(&todo); err != nil {
				return errors.WithStack(err)
			}

		}

		return nil
	})

	grift.Desc("deletetodoes", "Deletes all admin@wawand.co todoes")
	grift.Add("deletetodoes", func(c *grift.Context) error {

		// Get admin account...
		user := models.User{}
		if err := models.DB.Eager().Where("email = ?", "admin@wawand.co").First(&user); err != nil {
			errors.WithStack(err)
		}

		// Proceed deleting...
		if err := models.DB.Destroy(&user.Todoes); err != nil {
			errors.WithStack(err)
		}

		return nil

	})

})
