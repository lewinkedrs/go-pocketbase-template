package main

import (
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func loginHandler(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	//get credentials from request
	email := e.Request.FormValue("email")
	password := e.Request.FormValue("password")

	user, err := app.FindAuthRecordByEmail("users", email)
	if err != nil {
		return e.HTML(200, "User does not exist")
	}

	//Check if credentials are correct and return token
	isAuthenticated := user.ValidatePassword(password)
	if isAuthenticated {
		token, err := user.NewAuthToken()
		if err != nil {
			return e.HTML(200, "Error generating token, please try again.")
		}
		token_cookie := http.Cookie{
			Name:     "explore_token",
			Value:    token,
			Secure:   true,
			HttpOnly: true,
			Expires:  time.Now().Add(30 * time.Minute),
		}
		e.SetCookie(&token_cookie)
		id_cookie := http.Cookie{
			Name:     "id",
			Value:    user.Id,
			Secure:   true,
			HttpOnly: true,
			Expires:  time.Now().Add(30 * time.Minute),
		}
		e.SetCookie(&id_cookie)
		e.Response.Header().Set("HX-Redirect", "/")
	}

	loginError := "Login Failed, Please try again"
	return e.HTML(200, loginError)
}

func logoutHandler(e *core.RequestEvent) error {
	//clear cookies to logout
	token_cookie := new(http.Cookie)
	token_cookie.Name = "explore_token"
	token_cookie.Value = ""
	token_cookie.Expires = time.Now()
	e.SetCookie(token_cookie)

	id_cookie := new(http.Cookie)
	id_cookie.Name = "id"
	id_cookie.Value = ""
	id_cookie.Expires = time.Now()
	e.SetCookie(id_cookie)
	return e.Redirect(http.StatusMovedPermanently, "/")
}

func registerHandler(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	//extract registration data from request
	email := e.Request.FormValue("email")
	password := e.Request.FormValue("password")

	collection, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	//create auth record in users collection
	record := core.NewRecord(collection)
	record.SetEmail(email)
	record.SetPassword(password)

	err = app.Save(record)
	if err != nil {
		return err
	}

	return e.Redirect(302, "/")
}

func validateToken(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	//get token from request and validate it.
	token_cookie, err := e.Request.Cookie("explore_token")
	if err != nil {
		return err
	}
	_, err = app.FindAuthRecordByToken(token_cookie.Value, core.TokenTypeAuth)
	if err != nil {
		return err
	}
	return nil
}
