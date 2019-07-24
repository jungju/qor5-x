package setup

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/sunfmin/bran/examples/e17_hello_lazy_portals_and_reload"

	"github.com/sunfmin/bran/examples/e16_hello_vuetify_simple_components"

	"github.com/sunfmin/bran/examples/e15_hello_vuetify_navigation_drawer"

	"github.com/sunfmin/bran/examples/e13_hello_vuetify_list"
	"github.com/sunfmin/bran/examples/e14_hello_vuetify_menu"

	"github.com/sunfmin/bran/examples/e12_hello_vuetify_grid"

	"github.com/sunfmin/bran/codehighlight"
	"github.com/sunfmin/bran/core"
	branoverlay "github.com/sunfmin/bran/overlay"
	"github.com/sunfmin/bran/vuetify"

	"github.com/gobuffalo/packr"
	"github.com/sunfmin/bran"
	"github.com/sunfmin/bran/examples/e01_hello_button"
	"github.com/sunfmin/bran/examples/e02_hello_material_button"
	"github.com/sunfmin/bran/examples/e03_hello_card"
	"github.com/sunfmin/bran/examples/e04_hello_material_grid"
	"github.com/sunfmin/bran/examples/e06_hello_drawer"
	"github.com/sunfmin/bran/examples/e07_hello_lazy_portal_in_drawer"
	"github.com/sunfmin/bran/examples/e08_hello_popover"
	"github.com/sunfmin/bran/examples/e09_hello_dialog"
	"github.com/sunfmin/bran/examples/e10_hello_vuetify_autocomplete"
	"github.com/sunfmin/bran/examples/e11_hello_vuetify_text_field"
	m "github.com/sunfmin/bran/material"
	"github.com/sunfmin/bran/ui"
	"github.com/theplant/appkit/contexts"
	"github.com/theplant/appkit/server"
	. "github.com/theplant/htmlgo"
)

type pageItem struct {
	url        string
	renderFunc ui.PageFunc
	vuetify    bool
}

func (p pageItem) Title() string {
	segs := strings.Split(p.url, "_")
	segs[1] = strings.Title(segs[1])
	return fmt.Sprintf("%s: %s", strings.ToUpper(segs[0]), strings.Join(segs[1:], " "))
}

func exampleLinks(prefix string, pages []pageItem) (comp HTMLComponent) {
	var links []HTMLComponent
	for _, p := range pages {
		links = append(links,
			Li(
				A().Href(fmt.Sprintf("%s/%s", prefix, p.url)).Text(p.Title()),
			),
		)
	}
	comp = Ul(links...)
	return
}

var exampleBox = packr.NewBox("../")

func layout(in ui.PageFunc, pages []pageItem, prefix string, cp pageItem) (out ui.PageFunc) {
	return func(ctx *ui.EventContext) (pr ui.PageResponse, err error) {

		tailScript := `<script src='/assets/main.js'></script>`
		if len(os.Getenv("DEV")) > 0 {
			fmt.Println("Using Dev environment, make sure you did: yarn start")
			tailScript = `<script src='http://localhost:3100/app.js'></script>`
		}

		ctx.Injector.Title(cp.Title())
		ctx.Injector.PutHeadHTML(`
			<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto+Mono" async>
			<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500" async>
			<link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons" async>
			<link rel="stylesheet" href="/assets/main.css">
			<script src='/assets/vue.js'></script>
			<script src='/assets/codehighlight.js'></script>
		`)
		if cp.vuetify {
			ctx.Injector.PutHeadHTML(`
				<link rel="stylesheet" href="/assets/vuetify.css">
			`)

			tailScript := `<script src='/assets/vuetify.js'></script>`
			if len(os.Getenv("DEV")) > 0 {
				tailScript = `<script src='http://localhost:3080/app.js'></script>`
			}
			ctx.Injector.PutHeadHTML(tailScript)
		}

		ctx.Injector.PutTailHTML(tailScript)

		var innerPr ui.PageResponse
		innerPr, err = in(ctx)
		if err != nil {
			panic(err)
		}

		demo := innerPr.Schema

		var dacComps = []HTMLComponent{demo.(HTMLComponent)}

		var code string
		code, err = exampleBox.FindString(cp.url + "/page.go")
		if err != nil {
			return
		}
		if len(code) > 0 {
			dacComps = append(dacComps,
				Div(
					codehighlight.Code(code).Language("go"),
				).Class("exampleCode"),
			)
		}
		ctx.Injector.PutHeadHTML(`
		<style>
			body {
				margin: 0;
			}
			pre {
				padding: 0;
				margin: 0;
			}
			.exampleCode {
				margin-top: 20px;
			}
			[v-cloak] {
				display: none;
			}
		</style>
		`)

		pr.Schema = m.Grid(
			m.Cell(exampleLinks(prefix, pages)).Span(3, m.ScreenAll),
			m.Cell(dacComps...).Span(9, m.ScreenAll),
		)

		pr.State = innerPr.State

		return
	}
}

func home(prefix string, pages []pageItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/e01_hello_button", 302)
		return
	}
}

func Setup(prefix string) http.Handler {
	ub := bran.New()

	mux := http.NewServeMux()

	mux.Handle("/assets/main.js",
		ub.PacksHandler("text/javascript",
			branoverlay.JSComponentsPack(),
			core.JSComponentsPack(),
		),
	)

	mux.Handle("/assets/vue.js",
		ub.PacksHandler("text/javascript",
			core.JSVueComponentsPack(),
		),
	)

	mux.Handle("/assets/main.css",
		ub.PacksHandler("text/css",
			codehighlight.CSSComponentsPack(),
			branoverlay.CSSComponentsPack(),
			m.CSSComponentsPack(),
		),
	)

	mux.Handle("/assets/codehighlight.js",
		ub.PacksHandler("text/javascript",
			codehighlight.JSComponentsPack(),
		),
	)

	mux.Handle("/assets/vuetify.js",
		ub.PacksHandler("text/javascript",
			vuetify.JSComponentsPack(),
		),
	)

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	mux.Handle("/assets/vuetify.css",
		ub.PacksHandler("text/css",
			vuetify.CSSComponentsPack(),
		),
	)

	var pages = []pageItem{
		{
			url:        "e01_hello_button",
			renderFunc: e01_hello_button.HelloButton,
		},
		{
			url:        "e02_hello_material_button",
			renderFunc: e02_hello_material_button.HelloButton,
		},
		{
			url:        "e03_hello_card",
			renderFunc: e03_hello_card.HelloCard,
		},
		{
			url:        "e04_hello_material_grid",
			renderFunc: e04_hello_material_grid.HelloGrid,
		},
		//{
		//	url:        "e05_hello_customized_component",
		//	renderFunc: e05_hello_customized_component.HelloCustomziedComponent,
		//},
		{
			url:        "e06_hello_drawer",
			renderFunc: e06_hello_drawer.HelloDrawer,
		},
		{
			url:        "e07_hello_lazy_portal_in_drawer",
			renderFunc: e07_hello_lazy_portal_in_drawer.HelloLazyLoaderInDrawer,
		},
		{
			url:        "e08_hello_popover",
			renderFunc: e08_hello_popover.HelloPopover,
		},
		{
			url:        "e09_hello_dialog",
			renderFunc: e09_hello_dialog.HelloDialog,
		},
		{
			url:        "e10_hello_vuetify_autocomplete",
			renderFunc: e10_hello_vuetify_autocomplete.HelloVuetifyAutocomplete,
			vuetify:    true,
		},
		{
			url:        "e11_hello_vuetify_text_field",
			renderFunc: e11_hello_vuetify_text_field.HelloVuetifyTextField,
			vuetify:    true,
		},
		{
			url:        "e12_hello_vuetify_grid",
			renderFunc: e12_hello_vuetify_grid.HelloVuetifyGrid,
			vuetify:    true,
		},
		{
			url:        "e13_hello_vuetify_list",
			renderFunc: e13_hello_vuetify_list.HelloVuetifyList,
			vuetify:    true,
		},
		{
			url:        "e14_hello_vuetify_menu",
			renderFunc: e14_hello_vuetify_menu.HelloVuetifyMenu,
			vuetify:    true,
		},
		{
			url:        "e15_hello_vuetify_navigation_drawer",
			renderFunc: e15_hello_vuetify_navigation_drawer.HelloVuetifyNavigationDrawer,
			vuetify:    true,
		},
		{
			url:        "e16_hello_vuetify_simple_components",
			renderFunc: e16_hello_vuetify_simple_components.HelloVuetifySimpleComponents,
			vuetify:    true,
		},
		{
			url:        "e17_hello_lazy_portals_and_reload",
			renderFunc: e17_hello_lazy_portals_and_reload.HelloLazyPortalsAndReload,
			vuetify:    true,
		},
	}

	// l := log.Default()
	mw := server.Compose(
		// server.LogRequest,
		// log.WithLogger(l),
		contexts.WithHTTPStatus,
	)

	for _, p := range pages {
		mux.Handle(
			fmt.Sprintf("/%s", p.url),
			mw(ub.Page(layout(p.renderFunc, pages, prefix, p))),
		)
	}

	mux.Handle("/", home(prefix, pages))
	return mux
}
