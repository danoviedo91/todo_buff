package actions

import (
	"fmt"
	"log"
	"strconv"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// List is used to parse index.html the first time user enters the website
func List(c buffalo.Context) error {

	// Connect to the database

	tx := c.Value("tx").(*pop.Connection)

	todo := []models.Todo{}

	if c.Session().Get("current_user_id") == nil {
		c.Set("user", models.User{})
		return c.Render(200, r.HTML("index.html"))
	}

	userID := c.Session().Get("current_user_id").(uuid.UUID).String()

	if err := tx.Where("user_id = ?", userID).All(&todo); err != nil {
		log.Fatal(err)
	}

	// Assign values to variables...

	if c.Param("status") == "completed" {

		for i := 0; i < len(todo); i++ {
			if todo[i].Completed != true {
				todo = append(todo[:i], todo[i+1:]...)
				i--
			}
		}

	} else if c.Param("status") == "incompleted" {

		for i := 0; i < len(todo); i++ {
			if todo[i].Completed != false {
				todo = append(todo[:i], todo[i+1:]...)
				i--
			}
		}

	}

	var defaultMsg string

	if len(todo) == 0 {

		filterStatus := c.Session().Get("filterStatus").(string)

		if filterStatus == "completed" {
			defaultMsg = "No completed tasks to show"
		} else if filterStatus == "incompleted" {
			defaultMsg = "No incompleted tasks to show"
		} else {
			defaultMsg = "No tasks to show"
		}

	}

	// Prepare to send data to template

	c.Set("todo", todo)
	c.Set("mainViewFlag", true)
	c.Set("defaultMsg", defaultMsg)

	return c.Render(200, r.HTML("index.html"))

}

// NewTodo renders the todo/new.html which contains an empty form to create a todo
func NewTodo(c buffalo.Context) error {

	todo := models.Todo{}

	c.Set("todo", todo)

	return c.Render(200, r.HTML("todo/new.html"))

}

// CreateTodo creates a new todo inside the database as a row
func CreateTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	todo := models.Todo{}

	if err := c.Bind(&todo); err != nil {
		return err
	}

	todo.UserID = c.Session().Get("current_user_id").(uuid.UUID)

	// Validation
	verrs, err := todo.ValidateCreate(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("errors", verrs.Errors)
		c.Set("todo", todo)

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

	// Validate current user owns the todo....

	todo := models.Todo{}

	userID := c.Session().Get("current_user_id").(uuid.UUID).String()

	if err := tx.Where("user_id = ?", userID).Find(&todo, c.Param("todo_id")); err != nil {
		return c.Error(404, errors.New(fmt.Sprintf("could not find todo with id = %v", c.Param("todo_id"))))
	}

	// Delete...

	if err := tx.Destroy(&todo); err != nil {
		log.Fatal(err)
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

	todo := models.Todo{}

	userID := c.Session().Get("current_user_id").(uuid.UUID).String()

	if err := tx.Where("user_id = ?", userID).Find(&todo, c.Param("todo_id")); err != nil {
		return c.Error(404, errors.New(fmt.Sprintf("could not find todo with id = %v", c.Param("todo_id"))))
	}

	c.Set("todo", todo)
	c.Set("todoCurrentDate", strconv.Itoa(todo.Deadline.Year())+"-"+todo.MonthFormatted()+"-"+todo.DayFormatted())

	return c.Render(200, r.HTML("todo/edit.html"))
}

// UpdateTodo binds the edited info and updates the todo
func UpdateTodo(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)

	// Validate current user owns the todo....

	todo := models.Todo{}

	userID := c.Session().Get("current_user_id").(uuid.UUID).String()

	if err := tx.Where("user_id = ?", userID).Find(&todo, c.Param("ID")); err != nil {
		return c.Error(404, errors.New(fmt.Sprintf("could not find todo with id = %v", c.Param("ID"))))
	}

	// Then bind the form info to the struct...

	if err := c.Bind(&todo); err != nil {
		log.Fatal(err)
	}

	todo.UserID = c.Session().Get("current_user_id").(uuid.UUID)

	// Validate non-empty title...
	verrs, err := todo.ValidateUpdate(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("errors", verrs.Errors)
		c.Set("todo", todo)
		c.Set("todoCurrentDate", strconv.Itoa(todo.Deadline.Year())+"-"+todo.MonthFormatted()+"-"+todo.DayFormatted())
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

	todo := models.Todo{}

	userID := c.Session().Get("current_user_id").(uuid.UUID).String()

	if err := tx.Where("user_id = ?", userID).Find(&todo, c.Param("todo_id")); err != nil {
		return c.Error(404, errors.New(fmt.Sprintf("could not find todo with id = %v", c.Param("todo_id"))))
	}

	c.Set("todo", todo)

	return c.Render(200, r.HTML("todo/show.html"))

}
