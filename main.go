package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

		//get user's places hypermedia response
		e.Router.GET("/my_places", func(e *core.RequestEvent) error {
			id, err := e.Request.Cookie("id")
			if err != nil {
				return e.String(http.StatusOK, "Please log in to save a location")
			}
			id_cookie := id.Value

			auth, err := e.Request.Cookie("explore_token")
			if err != nil {
				return e.String(http.StatusOK, "Please log in to save a location")
			}
			data := fmt.Sprintf(`{"user_id":"%s"}`, id_cookie)
			url := fmt.Sprintf(`http://127.0.0.1:8090/api/collections/favorite_locations/records?filter=(user_id='%s')`, id_cookie)
			req, err := http.NewRequest("GET", url, strings.NewReader(data))
			if err != nil {
				return e.String(http.StatusInternalServerError, err.Error())
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", auth.Value)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				e.String(http.StatusInternalServerError, err.Error())
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return e.String(http.StatusInternalServerError, err.Error())
			}

			target := map[string]interface{}{}
			err = json.Unmarshal(body, &target)
			if err != nil {
				return err
			}

			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/places.html",
			).Render(map[string]any{
				"places": target,
			})
			if err != nil {
				return apis.NewNotFoundError("", err)
			}
			return e.HTML(http.StatusOK, html)
		})

		//auth handlers
		e.Router.POST("/logout", logoutHandler)
		e.Router.POST("/register", registerHandler)
		e.Router.POST("/loginHandler", loginHandler)

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
