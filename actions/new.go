package actions

import (
	"log"
	"time"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
)

// NewTodo default implementation.
func NewTodo(c buffalo.Context) error {

	// Establish database connection
	db := models.DB

	// Initialize database-query variables

	allRecords, records := []models.Todo{}, []models.Todo{}

	err := db.All(&allRecords)

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

	if c.Param("status") == "completed" {
		// If /?completed=true
		filterStatus.Completed = true
	} else if c.Param("status") == "incompleted" {
		// If /?completed=false
		filterStatus.Incompleted = true
	}

	// Prepare to send data to template

	c.Set("pendingTasksNumber", pendingTasksNumber)
	c.Set("currentDateTime", time.Now())
	c.Set("filterStatus", filterStatus)

	return c.Render(200, r.HTML("new/todo.html"))
}
