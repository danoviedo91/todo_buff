package actions

import (
	"log"
	"time"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {

	// Establish database connection
	db := models.DB

	// Initialize database-query variables

	allRecords, records := []models.Todo{}, []models.Todo{}
	pendingTasksNumber := 0

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

	// If /?completed=true

	if c.Param("status") == "completed" {
		for _, row := range allRecords {
			if row.Completed == true {
				records = append(records, row)
			}
		}
		pendingTasksNumber = len(allRecords) - len(records)
		filterStatus.Completed = true
		// If /?completed=false
	} else if c.Param("status") == "incompleted" {
		for _, row := range allRecords {
			if row.Completed == false {
				records = append(records, row)
			}
		}
		pendingTasksNumber = len(records)
		filterStatus.Incompleted = true
		// If /
	} else {
		for _, row := range allRecords {
			if row.Completed == false {
				records = append(records, row)
			}
		}
		pendingTasksNumber = len(records)
		records = allRecords
	}

	// Catch if there are tasks to show or not...

	defaultMsgFlag := false

	if len(records) == 0 {
		defaultMsgFlag = true
	}

	// Prepare to send data to template

	c.Set("defaultMsgFlag", defaultMsgFlag)
	c.Set("pendingTasksNumber", pendingTasksNumber)
	c.Set("currentDateTime", time.Now())
	c.Set("tasksArray", records)
	c.Set("filterStatus", filterStatus)
	c.Set("mainViewFlag", true)

	return c.Render(200, r.HTML("index.html"))
}
