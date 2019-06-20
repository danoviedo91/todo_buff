package middleware

import (
	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// HeaderInfo has information common among all handlers
func HeaderInfo(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		tx := c.Value("tx").(*pop.Connection)

		numberOfPendingTodoes, err := tx.Where("completed = ?", false).Count(&models.Todo{})

		if err != nil {
			return errors.WithStack(err)
		}

		c.Session().Set("filterStatus", c.Param("status"))
		c.Session().Set("numberOfPendingTodoes", numberOfPendingTodoes)

		return next(c)
	}
}
