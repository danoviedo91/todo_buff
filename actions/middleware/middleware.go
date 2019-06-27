package middleware

import (
	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// HeaderInfo has information common among all handlers
func HeaderInfo(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		tx := c.Value("tx").(*pop.Connection)
		var numberOfPendingTodoes int

		if c.Session().Get("current_user_id") != nil {

			var err error

			userID := c.Session().Get("current_user_id").(uuid.UUID).String()
			numberOfPendingTodoes, err = tx.Where("completed = ?", false).Where("user_id = ?", userID).Count(&models.Todo{})

			if err != nil {
				return errors.WithStack(err)
			}

		} else {

			numberOfPendingTodoes = 0

		}

		c.Session().Set("filterStatus", c.Param("status"))
		c.Session().Set("numberOfPendingTodoes", numberOfPendingTodoes)

		return next(c)
	}
}
