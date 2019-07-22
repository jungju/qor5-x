package e06_hello_drawer

import (
	"fmt"
	"math/rand"
	"time"

	bo "github.com/sunfmin/bran/overlay"
	"github.com/sunfmin/bran/ui"
	. "github.com/theplant/htmlgo"
)

type mystate struct {
	Name      string
	NameError string
	Group     string
}

func randStr(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s: %d", prefix, rand.Int31n(100))
}

func HelloDrawer(ctx *ui.EventContext) (pr ui.PageResponse, err error) {
	s := ctx.StateOrInit(&mystate{}).(*mystate)

	pr.Schema = Div(
		H1(s.Name),
		bo.Drawer(
			ui.LazyLoader(ctx.Hub, "form", form, "param1").LoadWhenParentVisible(),
		).TriggerElement(
			A().Text("Edit").Href("#"),
		).Width(500),
		ui.Bind(Input("").Type("text").Value(s.Group)).FieldName("Group"),
		ui.Bind(Button("Check value")).OnClick(ctx.Hub, "update", update),
	)
	return
}

func update(ctx *ui.EventContext) (r ui.EventResponse, err error) {
	r.Reload = true
	return
}

func form(ctx *ui.EventContext) (r ui.EventResponse, err error) {
	s := ctx.State.(*mystate)
	r.Schema = Div(
		Button("Close").Attr("@click", "parent.close"),
		ui.Bind(Input("").Type("text").Value(s.Name)).FieldName("Name"),
		Label(s.NameError).Style("color:red"),
		ui.Bind(Button("Update")).OnClick(ctx.Hub, "updateForm", updateForm),
	)
	return
}

func updateForm(ctx *ui.EventContext) (r ui.EventResponse, err error) {
	s := ctx.State.(*mystate)
	if len(s.Name) < 10 {
		s.NameError = "name is too short"
		s.Name = ""
		r, err = form(ctx)
	} else {
		s.NameError = ""
		r.Reload = true
	}
	return
}
