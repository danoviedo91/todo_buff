package actions

import (
	"fmt"
	"testing"
	"time"

	"github.com/danoviedo91/todo_buff/models"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {

	action, err := suite.NewActionWithFixtures(App(), packr.New("fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}

// All custom test functions start from this point

func (as *ActionSuite) Test_Empty_List_Todoes() {

	// -------- Logged out... -------- //

	res := as.HTML("/").Get()
	body := res.Body.String()

	as.Contains(body, `<div class="text-center wwc-mb-50">Welcome! Please log-in or register to continue</div>`)
	as.Equal(200, res.Code)

	// -------- Logged in... -------- //

	as.LoadFixture("sample user")

	user := &models.User{}
	as.NoError(as.DB.First(user))

	as.Session.Set("current_user_id", user.ID)

	res = as.HTML("/").Get()
	body = res.Body.String()

	as.Contains(body, fmt.Sprintf(`<td colspan="3" id="wwc-notasks-msg">No tasks to show</td>`))

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_List_All_Todoes() {

	// -------- Logged out... -------- //

	res := as.HTML("/").Get()
	body := res.Body.String()

	as.Contains(body, `<div class="text-center wwc-mb-50">Welcome! Please log-in or register to continue</div>`)
	as.Equal(200, res.Code)

	// -------- Logged in... -------- //

	as.LoadFixture("sample user and todoes")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	todo := &[]models.Todo{}
	as.NoError(as.DB.All(todo))

	res = as.HTML("/").Get()
	body = res.Body.String()

	for _, t := range *todo {
		as.Contains(body, fmt.Sprintf(`<a href="/show/%s/" class="wwc-task-title">%s</a>`, t.ID.String(), t.Title))
		as.Equal(t.UserID, user.ID)
	}

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Filtered_List_Todoes() { // ?status=complete

	// -------- Logged out... -------- //

	res := as.HTML("/?status=complete").Get()
	body := res.Body.String()

	as.Contains(body, `<div class="text-center wwc-mb-50">Welcome! Please log-in or register to continue</div>`)
	as.Equal(200, res.Code)

	// -------- Logged in... -------- //

	as.LoadFixture("sample user and todoes")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	todo := &[]models.Todo{}
	as.NoError(as.DB.All(todo))

	res = as.HTML("/?status=completed").Get()
	body = res.Body.String()

	for _, t := range *todo {

		if t.Completed == true {
			as.Contains(body, t.Title)
			as.Contains(body, t.ID.String())
		} else {
			as.NotContains(body, t.Title)
			as.NotContains(body, t.ID.String())
		}

		as.Equal(t.UserID, user.ID)

	}

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_New_Todo() {

	// -------- Logged out... -------- //

	res := as.HTML("/new").Get()
	body := res.Body.String()

	as.Contains(body, `<a href="/">Found</a>`) // From buffalo, this indicates redirection
	as.Equal(302, res.Code)
	as.Equal("/", res.Location())

	// -------- Logged in... -------- //

	as.LoadFixture("sample user")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	res = as.HTML("/new").Get()
	body = res.Body.String()

	as.Contains(body, "New Task")

	// Make sure form fields are empty
	as.Contains(body, `<input class=" form-control" id="todo-Title" name="Title" type="text" value="" />`)
	as.Contains(body, `<input class=" form-control" id="todo-Deadline" name="Deadline" type="date" value="" />`)
	as.Contains(body, `<textarea class=" form-control" id="todo-Description" name="Description" rows="4"></textarea>`)

	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Success_Create_Todo() {

	// -------- Logged out... -------- //

	todo := &models.Todo{}

	res := as.HTML("/create").Post(todo)
	body := res.Body.String()

	as.Contains(body, "") // From buffalo, this indicates redirection for Post Request
	as.Equal(302, res.Code)
	as.Equal("/", res.Location())

	// -------- Logged in... -------- //

	as.LoadFixture("sample user")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	todo = &models.Todo{
		Title:       "Hello World!",
		Description: "I am having success!",
		UserID:      user.ID,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Deadline:    time.Now(),
	}

	res = as.HTML("/create").Post(todo)

	err := as.DB.First(todo)
	as.NoError(err)
	as.NotZero(todo.ID)
	as.NotZero(todo.CreatedAt)
	as.Equal("Hello World!", todo.Title)
	as.Equal("I am having success!", todo.Description)
	as.Equal(todo.UserID, user.ID)
	as.Equal(todo.Completed, false)

	as.Equal(301, res.Code)
	as.Equal(fmt.Sprintf("/show/%s", todo.ID), res.Location())

}

func (as *ActionSuite) Test_Failed_Create_Todo() {

	// -------- Logged out... -------- //

	todo := &models.Todo{}

	res := as.HTML("/create").Post(todo)

	as.Equal(302, res.Code)
	as.Equal("/", res.Location())

	// -------- Logged in... -------- //

	as.LoadFixture("sample user")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	todo = &models.Todo{}

	res = as.HTML("/create").Post(todo)

	as.Error(as.DB.First(todo)) // Cannot create record, then error is expected
	as.Equal(422, res.Code)

}

func (as *ActionSuite) Test_Delete_Todo() {

	// -------- Logged out... -------- //

	res := as.HTML("/delete").Delete()
	as.Equal(404, res.Code)

	// -------- Logged in... -------- //

	as.LoadFixture("sample user and todoes")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	todo := &models.Todo{}

	as.NoError(as.DB.First(todo))

	as.Equal(todo.UserID, user.ID)

	as.Session.Set("filterStatus", "") // Necessary for redirection inside the handler

	res = as.HTML(fmt.Sprintf("/delete/%s", todo.ID.String())).Delete()

	as.Error(as.DB.Find(todo, todo.ID)) // Won't find the record, then error
	as.Equal(301, res.Code)
	as.Equal("/", res.Location())

}

func (as *ActionSuite) Test_Change_Status_Todo() {

	// -------- Logged out... -------- //

	todo := &models.Todo{}

	res := as.JSON("/change_status").Patch(todo)
	as.Equal(404, res.Code)

	// -------- Logged in... -------- //

	as.LoadFixture("sample user and todoes")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)
	as.Session.Set("filterStatus", "") // Necessary for redirection inside the handler

	// ------ Check when updates from false to true ------ //

	as.NoError(as.DB.Where("completed = ?", false).First(todo))
	todoID := todo.ID
	as.Equal(todo.UserID, user.ID)

	res = as.JSON(fmt.Sprintf("/change_status/%s", todo.ID.String())).Patch(todo)

	// After update, reload todo's data...
	as.NoError(as.DB.Find(todo, todoID))

	as.Equal(301, res.Code)
	as.Equal("/", res.Location())
	as.Equal(true, todo.Completed)

	// ------ Check when updates from true to false ------ //

	res = as.JSON(fmt.Sprintf("/change_status/%s", todo.ID.String())).Patch(todo)

	// After update, reload todo's data...
	as.NoError(as.DB.Find(todo, todoID))
	as.Equal(todo.UserID, user.ID)

	as.Equal(301, res.Code)
	as.Equal("/", res.Location())
	as.Equal(false, todo.Completed)

}

func (as *ActionSuite) Test_Show_Todo() {

	// -------- Logged out... -------- //

	todo := &models.Todo{}

	res := as.HTML("/show/samplestring").Get()
	body := res.Body.String()

	as.Contains(body, `<a href="/">Found</a>`) // From buffalo, this indicates redirection
	as.Equal(302, res.Code)
	as.Equal("/", res.Location())

	// -------- Logged in... Impostor testing -------- //

	as.LoadFixture("sample users and todoes")

	user := &models.User{}
	as.NoError(as.DB.Where("first_name = ?", "Bryan").First(user))
	as.Session.Set("current_user_id", user.ID)

	as.NoError(as.DB.Where("first_name = ?", "Daniel").First(user))
	as.DB.Where("user_id = ?", user.ID).First(todo)
	res = as.HTML(fmt.Sprintf("/show/%s", todo.ID.String())).Get()

	as.Equal(404, res.Code)

	// -------- Logged in... -------- //

	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	as.NoError(as.DB.Where("completed = ?", false).First(todo))
	as.Equal(todo.UserID, user.ID)

	res = as.HTML(fmt.Sprintf("/show/%s", todo.ID.String())).Get()
	body = res.Body.String()

	as.Contains(body, fmt.Sprintf(`<p><strong>Title:</strong>&nbsp;%s</p>`, todo.Title))
	as.Contains(body, fmt.Sprintf(`<p><strong>Description:</strong>&nbsp;%s</p>`, todo.Description))
	as.Contains(body, "No") // From todo.Completed : false
	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Edit_Todo() {

	// -------- Logged out... -------- //

	todo := &models.Todo{}

	res := as.HTML("/edit/samplestring").Get()
	body := res.Body.String()

	as.Contains(body, `<a href="/">Found</a>`) // From buffalo, this indicates redirection
	as.Equal(302, res.Code)
	as.Equal("/", res.Location())

	// -------- Logged in... -------- //

	as.LoadFixture("sample user and todoes")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	as.NoError(as.DB.First(todo))
	as.Equal(todo.UserID, user.ID)

	res = as.HTML(fmt.Sprintf("/edit/%s", todo.ID)).Get()
	body = res.Body.String()

	as.Contains(body, fmt.Sprintf(`<input class=" form-control" id="todo-Title" name="Title" type="text" value="%s" />`, todo.Title))
	as.Contains(body, fmt.Sprintf(`<textarea class=" form-control" id="todo-Description" name="Description" rows="4">%s</textarea>`, todo.Description))
	as.Contains(body, fmt.Sprintf(`<input class=" form-control" id="todo-Deadline" name="Deadline" type="date" value="%v-%v-%v" />`, todo.Deadline.Year(), todo.MonthFormatted(), todo.DayFormatted()))
	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Success_Update_Todo() {

	// -------- Logged out... -------- //

	todo := &models.Todo{}

	res := as.HTML("/update/samplestring").Put(todo)
	as.Equal(404, res.Code)

	// -------- Logged in... -------- //

	as.LoadFixture("sample user and todoes")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	as.NoError(as.DB.First(todo))
	todoID := todo.ID
	as.Equal(todo.UserID, user.ID)

	// todo - updated struct

	layout := "2006/01/02"
	timeString := "2019/03/03"
	timeParsed, err := time.Parse(layout, timeString)

	as.NoError(err)

	todo = &models.Todo{
		ID:          todo.ID,
		UserID:      user.ID,
		Title:       "Title Updated",
		Description: "Description Updated",
		Deadline:    timeParsed,
	}

	// Send the request with updated todo

	res = as.HTML("/update").Put(todo)

	// Verification

	as.NoError(as.DB.Find(todo, todoID))

	as.Equal("Title Updated", todo.Title)
	as.Equal("Description Updated", todo.Description)
	as.Equal(timeParsed, todo.Deadline.UTC())
	as.Equal(fmt.Sprintf("/show/%s", todo.ID), res.Location())
	as.Equal(301, res.Code)

}

func (as *ActionSuite) Test_Failed_Update_Todo() {

	// -------- Logged out... -------- //

	todo := &models.Todo{}

	res := as.HTML("/update/samplestring").Put(todo)
	as.Equal(404, res.Code)

	// -------- Logged in... -------- //

	as.LoadFixture("sample user and todoes")

	user := &models.User{}
	as.NoError(as.DB.First(user))
	as.Session.Set("current_user_id", user.ID)

	as.NoError(as.DB.First(todo))
	as.Equal(todo.UserID, user.ID)

	// todo - updated struct

	layout := "2006/01/02"
	timeString := "2019/02/02"
	timeParsed, err := time.Parse(layout, timeString)

	as.NoError(err)

	todo = &models.Todo{
		ID:          todo.ID,
		UserID:      user.ID,
		Description: "Description Updated",
		Deadline:    timeParsed,
	}

	// Send the request with invalid-updated todo

	res = as.HTML("/update").Put(todo)

	as.Equal(422, res.Code)

}
