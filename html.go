package main

import (
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func serveIndex(e *core.RequestEvent, registry *template.Registry) error {
	var isLoggedIn bool
	isLoggedIn = false
	cookie, _ := e.Request.Cookie("explore_token")
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
	return e.HTML(http.StatusOK, html)
}

func serveLogin(e *core.RequestEvent, registry *template.Registry) error {
	html, err := registry.LoadFiles(
		"views/layout.html",
		"views/login.html",
	).Render(map[string]any{
		"hello": "hello",
	})
	if err != nil {
		return apis.NewNotFoundError("", err)
	}
	return e.HTML(http.StatusOK, html)
}

func serveRegister(e *core.RequestEvent, registry *template.Registry) error {
	html, err := registry.LoadFiles(
		"views/layout.html",
		"views/register.html",
	).Render(map[string]any{
		"hello": "hello",
	})
	if err != nil {
		return apis.NewNotFoundError("", err)
	}
	return e.HTML(http.StatusOK, html)
}

func servePlaces(e *core.RequestEvent, registry *template.Registry, app *pocketbase.PocketBase) error {
	//validate token
	err := validateToken(e, app)
	if err != nil {
		return e.String(http.StatusOK, "Looks like you do not have a valid token, please login to save a location.")
	}

	//Extract user_id from request
	id, err := e.Request.Cookie("id")
	if err != nil {
		return e.String(http.StatusOK, "Please log in to save a location")
	}
	id_cookie := id.Value

	//Find all places for user.
	records, err := app.FindAllRecords("favorite_locations", dbx.NewExp("LOWER(`user_id`) = LOWER({:user_id})", dbx.Params{"user_id": id_cookie}))
	if err != nil {
		return err
	}

	var locationList []string
	for _, record := range records {
		locationName := record.GetString("location")
		if locationName != "" {
			locationList = append(locationList, locationName)
		}
	}

	//render template
	html, err := registry.LoadFiles(
		"views/layout.html",
		"views/places.html",
	).Render(map[string]any{
		"places": locationList,
	})
	if err != nil {
		return apis.NewNotFoundError("", err)
	}
	return e.HTML(http.StatusOK, html)
}
