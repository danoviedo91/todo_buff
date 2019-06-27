package models

import (
	"time"
)

// func Test_Todo(t *testing.T) {
// }

func (ms *ModelSuite) Test_Todo_Methods() {

	layout := "2006/01/02"
	timeString := "2019/02/02"
	timeParsed, err := time.Parse(layout, timeString)

	ms.NoError(err)

	todo := Todo{
		Deadline: timeParsed,
	}

	ms.Equal("02", todo.MonthFormatted())
	ms.Equal("02", todo.DayFormatted())
}
