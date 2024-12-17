package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// serve static files
		e.Router.GET("/static/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		// this is safe to be used by multiple goroutines
		// (it acts as store for the parsed templates)
		registry := template.NewRegistry()

		//hypermedia handlers
		//Serve up index
		e.Router.GET("/", func(e *core.RequestEvent) error {
			return serveIndex(e, registry)
		})
		//login page
		e.Router.GET("/login", func(e *core.RequestEvent) error {
			return serveLogin(e, registry)
		})
		//signup page
		e.Router.GET("/signup", func(e *core.RequestEvent) error {
			return serveRegister(e, registry)
		})
		//places hypermedia response
		e.Router.GET("/my_places", func(e *core.RequestEvent) error {
			return servePlaces(e, registry, app)
		})

		//auth handlers
		e.Router.POST("/logout", logoutHandler)
		//e.Router.POST("/register", registerHandler)
		e.Router.POST("/register", func(e *core.RequestEvent) error {
			return registerHandler(e, app)
		})
		e.Router.POST("/loginHandler", func(e *core.RequestEvent) error {
			return loginHandler(e, app)
		})

		//crud routes
		e.Router.POST("/save_location", func(e *core.RequestEvent) error {
			return saveLocation(e, app)
		})

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
