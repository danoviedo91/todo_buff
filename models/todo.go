package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Todo is the basic struct container of all the tasks info
type Todo struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Deadline    time.Time `json:"deadline" db:"deadline"`
	Completed   bool      `json:"completed" db:"completed"`
}

// String is not required by pop and may be deleted
func (t Todo) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Todoes is not required by pop and may be deleted
type Todoes []Todo

// String is not required by pop and may be deleted
func (t Todoes) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Todo) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Todo) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Title, Name: "Title", Message: "Todo title cannot be blank"},
		//&validators.TimeIsPresent{Field: t.Deadline, Name: "Deadline"},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Todo) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Title, Name: "Title", Message: "Todo title cannot be blank"},
		//&validators.TimeIsPresent{Field: t.Deadline, Name: "Deadline"},
	), nil
}

// ---- TODO FRONT END EDIT METHODS ---- //

// MonthFormatted converts type Month to type Int
func (t Todo) MonthFormatted() string {
	monthInt := int(t.Deadline.Month())
	return fmt.Sprintf("%02d", monthInt)
}

// DayFormatted converts type Month to type Int
func (t Todo) DayFormatted() string {
	return fmt.Sprintf("%02d", t.Deadline.Day())
}
