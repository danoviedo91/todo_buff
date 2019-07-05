package actions

import (
	"fmt"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
			"timeNow": func() string {

				t := time.Now().Format("Monday 2, January 2006")

				return t

			},
			"hrefCancelBtn": func(c buffalo.Context) string {

				filterStatus := c.Session().Get("filterStatus")
				page := c.Session().Get("page")
				var hrefCancelBtn string

				switch filterStatus {
				case "completed":
					hrefCancelBtn = fmt.Sprintf("/?page=%v&status=%v", page, "completed")
				case "incompleted":
					hrefCancelBtn = fmt.Sprintf("/?page=%v&status=%v", page, "incompleted")
				default:
					hrefCancelBtn = fmt.Sprintf("/?page=%v", page)
				}

				return hrefCancelBtn

			},
		},
	})
}
