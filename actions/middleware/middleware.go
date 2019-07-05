package middleware

import (
	"strconv"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// HeaderInfo has information common among all handlers
func HeaderInfo(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		var numberOfPendingTodoes int
		var err error

		tx := c.Value("tx").(*pop.Connection)

		if c.Value("current_user") != nil {

			userID := c.Value("current_user").(*models.User).ID
			isAdmin := c.Value("current_user").(*models.User).IsAdmin

			if isAdmin {

				numberOfPendingTodoes, err = tx.Where("completed = ?", false).Count(&models.Todo{})

				if err != nil {
					return errors.WithStack(err)
				}

			}

			if !isAdmin {

				numberOfPendingTodoes, err = tx.Where("completed = ?", false).Where("user_id = ?", userID).Count(&models.Todo{})

				if err != nil {
					return errors.WithStack(err)
				}

			}

		}

		page := 1

		if c.Param("page") != "" {
			page, err = strconv.Atoi(c.Param("page"))
			if err != nil {
				errors.WithStack(err)
			}
		}

		c.Session().Set("page", page)
		c.Session().Set("filterStatus", c.Param("status"))
		c.Session().Set("numberOfPendingTodoes", numberOfPendingTodoes)

		return next(c)
	}
}

// AuthorizeAdmin verifies the current user is admin or not, and if not, redirect to main dashboard
func AuthorizeAdmin(next buffalo.Handler) buffalo.Handler {

	return func(c buffalo.Context) error {

		isAdmin := c.Value("current_user").(*models.User).IsAdmin

		if !isAdmin {

			return c.Redirect(301, "/")

		}

		return next(c)

	}
}
