package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serve static files
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		// this is safe to be used by multiple goroutines
		// (it acts as store for the parsed templates)
		registry := template.NewRegistry()

		e.Router.GET("/hello/:name", func(c echo.Context) error {
			name := c.PathParam("name")

			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/hello.html",
			).Render(map[string]any{
				"name": name,
			})
			if err != nil {
				// or redirect to a dedicated 404 HTML page
				return apis.NewNotFoundError("", err)
			}

			return c.HTML(http.StatusOK, html)
		})

		e.Router.GET("/", func(c echo.Context) error {
			var isLoggedIn bool
			isLoggedIn = false
			cookie, _ := c.Cookie("explore_token")
			if cookie != nil {
				isLoggedIn = true
			}
			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/nav.html",
				"views/index.html",
			).Render(map[string]any{
				"isLoggedIn": isLoggedIn,
			})
			if err != nil {
				return apis.NewNotFoundError("", err)
			}
			return c.HTML(http.StatusOK, html)
		})

		e.Router.GET("/login", func(c echo.Context) error {
			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/login.html",
			).Render(map[string]any{
				"hello": "hello",
			})
			if err != nil {
				return apis.NewNotFoundError("", err)
			}
			return c.HTML(http.StatusOK, html)
		})

		e.Router.POST("/logout", func(c echo.Context) error {
			token_cookie := new(http.Cookie)
			token_cookie.Name = "explore_token"
			token_cookie.Value = ""
			token_cookie.Expires = time.Now()
			c.SetCookie(token_cookie)
			id_cookie := new(http.Cookie)
			id_cookie.Name = "id"
			id_cookie.Value = ""
			id_cookie.Expires = time.Now()
			c.SetCookie(id_cookie)
			return c.Redirect(http.StatusMovedPermanently, "/")
		})

		e.Router.GET("/signup", func(c echo.Context) error {
			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/register.html",
			).Render(map[string]any{
				"hello": "hello",
			})
			if err != nil {
				return apis.NewNotFoundError("", err)
			}
			return c.HTML(http.StatusOK, html)
		})

		e.Router.GET("/my_places", func(c echo.Context) error {
			id, err := c.Cookie("id")
			if err != nil {
				return c.String(http.StatusOK, "Please log in to save a location")
			}
			id_cookie := id.Value

			auth, err := c.Cookie("explore_token")
			if err != nil {
				return c.String(http.StatusOK, "Please log in to save a location")
			}
			url := fmt.Sprintf(`http://127.0.0.1:8090/api/collections/favorite_locations/records?filter=(user_id='%s')`, id_cookie)
			req, err := http.NewRequest("GET", url, strings.NewReader(""))
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", auth.Value)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
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
			return c.HTML(http.StatusOK, html)
		})

		e.Router.POST("/register", func(c echo.Context) error {
			email := c.FormValue("email")
			password := c.FormValue("password")

			data := fmt.Sprintf(`{"email":"%s","password":"%s","passwordConfirm":"%s"}`, email, password, password)

			req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/users/records", strings.NewReader(data))
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			if resp.StatusCode == http.StatusOK {
				return c.Redirect(http.StatusMovedPermanently, "/")
			}

			return c.String(http.StatusBadRequest, string(body))
		})

		e.Router.POST("save_location", func(c echo.Context) error {
			location := c.FormValue("location")

			id, err := c.Cookie("id")
			if err != nil {
				return c.String(http.StatusOK, "Please log in to save a location")
			}
			id_cookie := id.Value

			auth, err := c.Cookie("explore_token")
			if err != nil {
				return c.String(http.StatusOK, "Please log in to save a location")
			}

			data := fmt.Sprintf(`{"location":"%s","user_id":"%s"}`, location, id_cookie)
			req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/favorite_locations/records", strings.NewReader(data))
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", auth.Value)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			if resp.StatusCode == http.StatusOK {
				c.String(http.StatusOK, "Successfully Saved Location")
			}

			return c.String(http.StatusBadRequest, string(body))
		})

		e.Router.POST("/loginHandler", func(c echo.Context) error {
			email := c.FormValue("email")
			password := c.FormValue("password")

			data := fmt.Sprintf(`{"identity":"%s","password":"%s","passwordConfirm":"%s"}`, email, password, password)

			req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/users/auth-with-password", strings.NewReader(data))
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			if resp.StatusCode == http.StatusOK {
				var authResponse struct {
					Token  string `json:"token"`
					Record struct {
						ID              string `json:"id"`
						CollectionID    string `json:"collectionId"`
						CollectionName  string `json:"collectionName"`
						Created         string `json:"created"`
						Updated         string `json:"updated"`
						Username        string `json:"username"`
						Email           string `json:"email"`
						Verified        bool   `json:"verified"`
						EmailVisibility bool   `json:"emailVisibility"`
						SomeCustomField string `json:"someCustomField"`
					} `json:"record"`
				}
				err = json.Unmarshal(body, &authResponse)
				if err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
				token_cookie := new(http.Cookie)
				token_cookie.Name = "explore_token"
				token_cookie.Value = authResponse.Token
				token_cookie.Expires = time.Now().Add(24 * time.Hour)
				c.SetCookie(token_cookie)
				id_cookie := new(http.Cookie)
				id_cookie.Name = "id"
				id_cookie.Value = authResponse.Record.ID
				id_cookie.Expires = time.Now().Add(24 * time.Hour)
				c.SetCookie(id_cookie)
				return c.Redirect(http.StatusMovedPermanently, "/")
			}

			return c.String(http.StatusBadRequest, string(body))
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
