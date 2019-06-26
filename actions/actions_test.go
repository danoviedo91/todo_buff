package actions

import (
	"fmt"
	"testing"
	"time"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite"
	"github.com/gofrs/uuid"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.New("Test_ActionSuite", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}

	suite.Run(t, as)
}

// All custom test functions start from this point

func (as *ActionSuite) Test_Todo_Empty_List() {

	res := as.HTML("/").Get()
	body := res.Body.String()

	as.Contains(body, fmt.Sprintf(`<td colspan="3" id="wwc-notasks-msg">No tasks to show</td>`))

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_List_Todo() {

	todo := []models.Todo{
		{
			Title: "Washing the dog",
			ID:    uuid.Must(uuid.NewV4()),
		},
		{
			Title: "Writing a book",
			ID:    uuid.Must(uuid.NewV4()),
		},
	}

	for _, t := range todo {
		err := as.DB.Create(&t)
		as.NoError(err)
	}

	res := as.HTML("/").Get()
	body := res.Body.String()

	for _, t := range todo {
		as.Contains(body, fmt.Sprintf(`<a href="/show/%s/" class="wwc-task-title">%s</a>`, t.ID.String(), t.Title))
	}

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Filtered_List_Todo() { // ?status=complete

	todo := []models.Todo{
		{
			Title:     "Washing the dog",
			ID:        uuid.Must(uuid.NewV4()),
			Completed: true,
		},
		{
			Title:     "Writing an essay",
			ID:        uuid.Must(uuid.NewV4()),
			Completed: true,
		},
		{
			Title:     "Writing a book",
			ID:        uuid.Must(uuid.NewV4()),
			Completed: false,
		},
		{
			Title:     "Washing the cat",
			Completed: false,
		},
	}

	for _, t := range todo {
		err := as.DB.Create(&t)
		as.NoError(err)
	}

	res := as.HTML("/?status=completed").Get()
	body := res.Body.String()

	for _, t := range todo {

		if t.Completed == true {
			as.Contains(body, t.Title)
			as.Contains(body, t.ID.String())
		} else {
			as.NotContains(body, t.Title)
			as.NotContains(body, t.ID.String())
		}

	}

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_New_Todo() {

	res := as.HTML("/new").Get()
	body := res.Body.String()

	as.Contains(body, "New Task")

	// Make sure form fields are empty
	as.Contains(body, `<input class=" form-control" id="todo-Title" name="Title" type="text" value="" />`)
	as.Contains(body, `<input class=" form-control" id="todo-Deadline" name="Deadline" type="date" value="" />`)
	as.Contains(body, `<textarea class=" form-control" id="todo-Description" name="Description" rows="4"></textarea>`)

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Success_Create_Todo() {

	todo := &models.Todo{
		Title:       "Hello World!",
		Description: "I am having success!",
	}

	res := as.HTML("/create").Post(todo)

	err := as.DB.First(todo)
	as.NoError(err)
	as.NotZero(todo.ID)
	as.NotZero(todo.CreatedAt)
	as.Equal("Hello World!", todo.Title)
	as.Equal("I am having success!", todo.Description)

	as.Equal(301, res.Code)
	as.Equal(fmt.Sprintf("/show/%s", todo.ID), res.Location())

}

func (as *ActionSuite) Test_Failed_Create_Todo() {

	todo := &models.Todo{}

	res := as.HTML("/create").Post(todo)

	err := as.DB.First(todo)
	as.Error(err)
	as.Equal(422, res.Code)

}

func (as *ActionSuite) Test_Delete_Todo() {

	todo := &models.Todo{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       "I'm being deleted",
		Description: "Some deleteable description",
	}

	as.Session.Set("filterStatus", "") // Necessary for redirection inside the handler

	err := as.DB.Create(todo)
	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/delete/%s", todo.ID.String())).Delete()

	as.Equal(301, res.Code)
	as.Equal("/", res.Location())

}

func (as *ActionSuite) Test_Change_Status_Todo() {

	todo := &models.Todo{
		ID:    uuid.Must(uuid.NewV4()),
		Title: "I'm being completed",
	}

	as.Session.Set("filterStatus", "") // Necessary for redirection inside the handler

	err := as.DB.Create(todo)
	as.NoError(err)

	// Check when updates from false to true

	res := as.JSON(fmt.Sprintf("/change_status/%s", todo.ID.String())).Patch(todo)

	err = as.DB.First(todo)
	as.NoError(err)

	as.Equal(301, res.Code)
	as.Equal("/", res.Location())
	as.Equal(true, todo.Completed)

	// Check when updates from true to false

	res = as.JSON(fmt.Sprintf("/change_status/%s", todo.ID.String())).Patch(todo)

	err = as.DB.First(todo)
	as.NoError(err)

	as.Equal(301, res.Code)
	as.Equal("/", res.Location())
	as.Equal(false, todo.Completed)

}

func (as *ActionSuite) Test_Show_Todo() {

	todo := &models.Todo{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       "I am being showed up!",
		Description: "A description to be showed up!",
		Completed:   false,
	}

	err := as.DB.Create(todo)
	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/show/%s", todo.ID.String())).Get()
	body := res.Body.String()

	as.Contains(body, fmt.Sprintf(`<p><strong>Title:</strong>&nbsp;%s</p>`, todo.Title))
	as.Contains(body, fmt.Sprintf(`<p><strong>Description:</strong>&nbsp;%s</p>`, todo.Description))
	as.Contains(body, "No") // From todo.Completed : false
	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Edit_Todo() {

	layout := "2006/01/02"
	timeString := "2019/02/02"
	timeParsed, err := time.Parse(layout, timeString)

	as.NoError(err)

	todo := &models.Todo{
		Title:       "I will edit this task",
		Description: "I will edit this description",
		Deadline:    timeParsed,
		Completed:   false,
	}

	err = as.DB.Create(todo)

	as.NoError(err)

	res := as.HTML(fmt.Sprintf("/edit/%s", todo.ID)).Get()
	body := res.Body.String()

	as.Contains(body, fmt.Sprintf(`<input class=" form-control" id="todo-Title" name="Title" type="text" value="%s" />`, todo.Title))
	as.Contains(body, fmt.Sprintf(`<textarea class=" form-control" id="todo-Description" name="Description" rows="4">%s</textarea>`, todo.Description))
	as.Contains(body, fmt.Sprintf(`<input class=" form-control" id="todo-Deadline" name="Deadline" type="date" value="%v-%v-%v" />`, todo.Deadline.Year(), todo.MonthFormatted(), todo.DayFormatted()))
	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Success_Update_Todo() {

	id := uuid.Must(uuid.NewV4())

	// Initial todo, create record in the db

	layout := "2006/01/02"
	timeString := "2019/02/02"
	timeParsed, err := time.Parse(layout, timeString)

	as.NoError(err)

	todo := &models.Todo{
		ID:          id,
		Title:       "I will edit this task",
		Description: "I will edit this description",
		Deadline:    timeParsed,
		Completed:   false,
	}

	err = as.DB.Create(todo)

	as.NoError(err)

	// todo - updated struct

	layout = "2006/01/02"
	timeString = "2019/03/03"
	timeParsed, err = time.Parse(layout, timeString)

	as.NoError(err)

	todo = &models.Todo{
		ID:          id,
		Title:       "Title Updated",
		Description: "Description Updated",
		Deadline:    timeParsed,
		Completed:   false,
	}

	// Send the request with updated todo

	res := as.HTML("/update").Put(todo)

	// Verification

	as.DB.Find(todo, id)

	as.Equal("Title Updated", todo.Title)
	as.Equal("Description Updated", todo.Description)
	as.Equal(timeParsed, todo.Deadline.UTC())
	as.Equal(fmt.Sprintf("/show/%s", todo.ID), res.Location())
	as.Equal(301, res.Code)

}

func (as *ActionSuite) Test_Failed_Update_Todo() {

	id := uuid.Must(uuid.NewV4())

	// Initial todo, create record in the db

	layout := "2006/01/02"
	timeString := "2019/02/02"
	timeParsed, err := time.Parse(layout, timeString)

	as.NoError(err)

	todo := &models.Todo{
		ID:          id,
		Title:       "I will edit this task",
		Description: "I will edit this description",
		Deadline:    timeParsed,
		Completed:   false,
	}

	err = as.DB.Create(todo)

	as.NoError(err)

	// Invalid todo - updated struct with NO title

	layout = "2006/01/02"
	timeString = "2019/03/03"
	timeParsed, err = time.Parse(layout, timeString)

	as.NoError(err)

	todo = &models.Todo{
		ID:          id,
		Description: "Description Updated",
		Deadline:    timeParsed,
		Completed:   false,
	}

	// Send the request with invalid-updated todo

	res := as.HTML("/update").Put(todo)

	as.Equal(422, res.Code)

}
