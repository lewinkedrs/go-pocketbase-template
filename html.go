package main

import (
	"net/http"

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
