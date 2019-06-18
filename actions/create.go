package actions

import (
	"fmt"
	"log"
	"time"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
)

// CreateTodo default implementation.
func CreateTodo(c buffalo.Context) error {

	// Establish database connection
	db := models.DB

	newTodo := models.Todo{}

	// Get form info
	if err := c.Bind(&newTodo); err != nil {
		return err
	}

	// Insert record into the database
	err := db.Create(&newTodo)

	if err != nil {
		fmt.Println(err)
	}

	// Initialize database-query variables

	allRecords, records := []models.Todo{}, []models.Todo{}

	err = db.All(&allRecords)

	if err != nil {
		log.Fatal(err)
	}

	filterStatus := struct {
		Incompleted bool
		Completed   bool
	}{
		false,
		false,
	}

	// Assign values to variables...

	for _, row := range allRecords {
		if row.Completed == false {
			records = append(records, row)
		}
	}
	pendingTasksNumber := len(records)
	records = allRecords

	// Prepare to send data to template

	c.Set("pendingTasksNumber", pendingTasksNumber)
	c.Set("currentDateTime", time.Now())
	c.Set("taskStruct", newTodo)
	c.Set("filterStatus", filterStatus)

	return c.Render(200, r.HTML("show/todo.html"))
}
