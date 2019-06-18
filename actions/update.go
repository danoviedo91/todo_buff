package actions

import (
	"fmt"
	"log"
	"time"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
)

// UpdateTodo default implementation.
func UpdateTodo(c buffalo.Context) error {

	// Establish database connection
	db := models.DB

	updateTodo := models.Todo{}

	err := c.Bind(&updateTodo)

	if err != nil {
		return err
	}

	err = db.Update(&updateTodo)

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
	c.Set("taskStruct", updateTodo)
	c.Set("filterStatus", filterStatus)

	return c.Render(200, r.HTML("show/todo.html"))
}
