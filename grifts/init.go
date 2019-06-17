package grifts

import (
	"github.com/danoviedo91/todo_buff/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
