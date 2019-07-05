package actions

import (
	"fmt"
	"log"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// List is used to parse index.html the first time user enters the website
func List(c buffalo.Context) error {

	// Render template immediately if not logged in...

	if c.Value("current_user") == nil {
		return c.Redirect(301, "/signin")
	}

	isAdmin := c.Value("current_user").(*models.User).IsAdmin

	// Retrieve user's todos
	tx := c.Value("tx").(*pop.Connection)

	q, todo, err := getTodoes(tx, c)

	if err != nil {
		return errors.WithStack(err)
	}

	var defaultMsg string

	if len(todo) == 0 {

		filterStatus := c.Session().Get("filterStatus").(string)

		switch filterStatus {

		case "completed":
			defaultMsg = "No completed tasks to show"
		case "incompleted":
			defaultMsg = "No incompleted tasks to show"
		default:
			defaultMsg = "No tasks to show"

		}

	}

	users := []models.User{}

	if isAdmin {
		if err := tx.All(&users); err != nil {
			return errors.WithStack(err)
		}
	}

	// Prepare to send data to template

	c.Set("user_id", c.Session().Get("current_user_id"))
	c.Set("users", users)
	c.Set("isAdmin", isAdmin)
	c.Set("todo", todo)
	c.Set("mainViewFlag", true)
	c.Set("pagination", q.Paginator)
	c.Set("defaultMsg", defaultMsg)

	return c.Render(200, r.HTML("index.html"))

}

// NewTodo renders the todo/new.html which contains an empty form to create a todo
func NewTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)

	todo := models.Todo{}
	users := []models.User{}
	// userOptions, userID, isAdmin - used for the select tag inside the front end form
	userOptions := make(map[string]string)
	userID := c.Value("current_user").(*models.User).ID
	isAdmin := c.Value("current_user").(*models.User).IsAdmin

	if err := tx.All(&users); err != nil {
		return errors.WithStack(err)
	}

	for _, user := range users {
		optionValue := user.ID.String()
		optionName := fmt.Sprintf("%v %v - %v", user.FirstName, user.LastName, user.Email)
		userOptions[optionName] = optionValue
	}

	// Setting context variables

	c.Set("isAdmin", isAdmin)
	c.Set("userID", userID)
	c.Set("userOptions", userOptions)
	c.Set("todo", todo)
	c.Set("context", c) // Used for HrefCancelBtn helper

	return c.Render(200, r.HTML("todo/new.html"))

}

// CreateTodo creates a new todo inside the database as a row
func CreateTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	todo := models.Todo{}

	if err := c.Bind(&todo); err != nil {
		return err
	}

	if c.Param("UserID") == "" {
		todo.UserID = c.Value("current_user").(*models.User).ID
	}

	todo.Completed = false

	// Validation
	verrs, err := todo.ValidateCreate(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {

		users := []models.User{}
		userOptions := make(map[string]string)
		userID := c.Value("current_user").(*models.User).ID
		isAdmin := c.Value("current_user").(*models.User).IsAdmin

		if err := tx.All(&users); err != nil {
			return errors.WithStack(err)
		}

		for _, user := range users {
			optionValue := user.ID.String()
			optionName := fmt.Sprintf("%v %v - %v", user.FirstName, user.LastName, user.Email)
			userOptions[optionName] = optionValue
		}

		c.Set("errors", verrs.Errors)
		c.Set("isAdmin", isAdmin)
		c.Set("userID", userID)
		c.Set("userOptions", userOptions)
		c.Set("todo", todo)
		c.Set("context", c) // Used for HrefCancelBtn helper

		return c.Render(422, r.HTML("todo/new.html"))
	}

	// proceed, no errors found
	// c.Flash().Add("success", "New post added successfully.")

	if err := tx.Create(&todo); err != nil {
		return errors.WithStack(err)
	}

	// Prepare to send data to template

	return c.Redirect(301, "/show/%v", todo.ID.String())

}

// DeleteTodo removes a todo from the database given an ID
func DeleteTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)

	todo, err := getTodo(tx, c)

	if err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Destroy(&todo); err != nil {
		return errors.WithStack(err)
	}

	// Redirect to "/" with filter, if applies

	var status string

	if c.Session().Get("filterStatus") != "" {

		status = "?status=" + c.Session().Get("filterStatus").(string)

	}

	return c.Redirect(301, "/%v", status)

}

// CompleteTodo completes a todo given an ID
func CompleteTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	todo := models.Todo{}

	if err := tx.Find(&todo, c.Param("todo_id")); err != nil {
		log.Fatal(err)
	}

	// Update record given the id

	todo.Completed = !todo.Completed

	if err := tx.Update(&todo); err != nil {
		log.Fatal(err)
	}

	// Redirect to "/" with filter, if applies

	status := c.Session().Get("filterStatus").(string)

	if status != "" {

		status = "?status=" + status

	}

	return c.Redirect(301, "/%v", status)

}

// EditTodo renders the todo/edit.html which contains a form with the todo's info
func EditTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	users := []models.User{}
	userOptions := make(map[string]string)
	isAdmin := c.Value("current_user").(*models.User).IsAdmin

	todo, err := getTodo(tx, c)

	if err := tx.All(&users); err != nil {
		return errors.WithStack(err)
	}

	for _, user := range users {
		optionValue := user.ID.String()
		optionName := fmt.Sprintf("%v %v - %v", user.FirstName, user.LastName, user.Email)
		userOptions[optionName] = optionValue
	}

	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("isAdmin", isAdmin)
	c.Set("userOptions", userOptions)
	c.Set("todo", todo)
	c.Set("todoCurrentDate", fmt.Sprintf("%v-%v-%v", todo.Deadline.Year(), todo.MonthFormatted(), todo.DayFormatted()))
	c.Set("context", c) // Used for HrefCancelBtn helper

	return c.Render(200, r.HTML("todo/edit.html"))
}

// UpdateTodo binds the edited info and updates the todo
func UpdateTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)

	todo, err := getTodo(tx, c)

	if err != nil {
		return errors.WithStack(err)
	}

	// Then bind the form info to the struct...

	if err := c.Bind(&todo); err != nil {
		log.Fatal(err)
	}

	if c.Param("UserID") == "" {
		todo.UserID = c.Value("current_user").(*models.User).ID
	}

	// Validate non-empty title...
	verrs, err := todo.ValidateUpdate(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {

		users := []models.User{}
		userOptions := make(map[string]string)
		userID := c.Value("current_user").(*models.User).ID
		isAdmin := c.Value("current_user").(*models.User).IsAdmin

		if err := tx.All(&users); err != nil {
			return errors.WithStack(err)
		}

		for _, user := range users {
			optionValue := user.ID.String()
			optionName := fmt.Sprintf("%v %v - %v", user.FirstName, user.LastName, user.Email)
			userOptions[optionName] = optionValue
		}

		c.Set("isAdmin", isAdmin)
		c.Set("userID", userID)
		c.Set("userOptions", userOptions)
		c.Set("errors", verrs.Errors)
		c.Set("todo", todo)
		c.Set("context", c) // Used for HrefCancelBtn helper
		c.Set("todoCurrentDate", fmt.Sprintf("%v-%v-%v", todo.Deadline.Year(), todo.MonthFormatted(), todo.DayFormatted()))

		return c.Render(422, r.HTML("todo/edit.html"))

	}

	// Update...

	if err := tx.Update(&todo); err != nil {
		log.Fatal(err)
	}

	return c.Redirect(301, "/show/%v", todo.ID.String())
}

// ShowTodo displays all the information of a todo
func ShowTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}

	todo, err := getTodo(tx, c)

	if err := tx.Find(&user, todo.UserID); err != nil {
		return errors.WithStack(err)
	}

	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("todo", todo)
	c.Set("owner", fmt.Sprintf("%v %v", user.FirstName, user.LastName))
	c.Set("context", c) // Used for HrefCancelBtn helper

	return c.Render(200, r.HTML("todo/show.html"))

}

// ============================ ----
// ---- In-action helper functions...
// ============================ ----

func getTodoes(tx *pop.Connection, c buffalo.Context) (*pop.Query, []models.Todo, error) {

	q := tx.PaginateFromParams(c.Params())
	q.Paginator.PerPage = 5
	q.Paginator.Offset = (q.Paginator.Page - 1) * q.Paginator.PerPage

	todo := []models.Todo{}

	userID := c.Value("current_user").(*models.User).ID
	isAdmin := c.Value("current_user").(*models.User).IsAdmin

	filter := c.Param("status")

	if isAdmin {

		switch filter {
		case "completed":
			if err := q.Where("completed = ?", true).All(&todo); err != nil {
				return nil, nil, errors.WithStack(err)
			}
		case "incompleted":
			if err := q.Where("completed = ?", false).All(&todo); err != nil {
				return nil, nil, errors.WithStack(err)
			}
		default:
			if err := q.All(&todo); err != nil {
				return nil, nil, errors.WithStack(err)
			}
		}

	}

	if !isAdmin {

		switch filter {
		case "completed":
			if err := q.Where("user_id = ?", userID).Where("completed = ?", true).All(&todo); err != nil {
				return nil, nil, errors.WithStack(err)
			}
		case "incompleted":
			if err := q.Where("user_id = ?", userID).Where("completed = ?", false).All(&todo); err != nil {
				return nil, nil, errors.WithStack(err)
			}
		default:
			if err := q.Where("user_id = ?", userID).All(&todo); err != nil {
				return nil, nil, errors.WithStack(err)
			}
		}

	}

	return q, todo, nil

}

func getTodo(tx *pop.Connection, c buffalo.Context) (models.Todo, error) {

	todo := models.Todo{}
	userID := c.Value("current_user").(*models.User).ID
	isAdmin := c.Value("current_user").(*models.User).IsAdmin

	if isAdmin {
		if err := tx.Find(&todo, c.Param("todo_id")); err != nil {
			return models.Todo{}, errors.WithStack(err)
		}
	}

	if !isAdmin {
		if err := tx.Where("user_id = ?", userID).Find(&todo, c.Param("todo_id")); err != nil {
			return models.Todo{}, c.Error(404, errors.New(fmt.Sprintf("could not find todo with id = %v", c.Param("todo_id"))))
		}
	}

	return todo, nil

}
