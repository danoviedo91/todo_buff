package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//User is a generated model from buffalo-auth, it serves as the base for username/password authentication.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email" db:"email"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`

	Todoes               []Todo `has_many:"todoes" db:"-"`
	Password             string `json:"-" db:"-"`
	PasswordConfirmation string `json:"-" db:"-"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	return tx.ValidateAndCreate(u)
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.RegexMatch{Field: u.Password, Name: "Password", Expr: "[^A-Za-z0-9]", Message: "Password must contain at least one special character"},
		&validators.RegexMatch{Field: u.Password, Name: "Password", Expr: ".*[A-Z].*", Message: "Password must contain at least one uppercase letter"},
		&validators.RegexMatch{Field: u.Password, Name: "Password", Expr: ".*[a-z].*", Message: "Password must contain at least one lowercase letter"},
		&validators.StringLengthInRange{Field: u.Password, Name: "Password", Min: 8, Message: "Password must be at least 8 characters long"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// GetTodoes ...
func (u User) GetTodoes(tx *pop.Connection, c buffalo.Context) ([]Todo, error) {

	todo := []Todo{}
	userID := c.Value("current_user").(*User).ID
	isAdmin := c.Value("current_user").(*User).IsAdmin

	if isAdmin {

		if err := tx.All(&todo); err != nil {
			return nil, errors.WithStack(err)
		}

	}

	if !isAdmin {

		if err := tx.Eager().Find(&u, userID); err != nil {
			return nil, errors.WithStack(err)
		}

		todo = u.Todoes

	}

	return todo, nil

}

// GetTodo ...
func (u User) GetTodo(tx *pop.Connection, c buffalo.Context) (Todo, error) {

	todo := Todo{}
	userID := c.Value("current_user").(*User).ID
	isAdmin := c.Value("current_user").(*User).IsAdmin

	if isAdmin {
		if err := tx.Find(&todo, c.Param("todo_id")); err != nil {
			return Todo{}, errors.WithStack(err)
		}
	}

	if !isAdmin {
		if err := tx.Where("user_id = ?", userID).Find(&todo, c.Param("todo_id")); err != nil {
			return Todo{}, c.Error(404, errors.New(fmt.Sprintf("could not find todo with id = %v", c.Param("todo_id"))))
		}
	}

	return todo, nil

}
