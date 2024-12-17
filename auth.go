package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func loginHandler(e *core.RequestEvent) error {
	email := e.Request.FormValue("email")
	password := e.Request.FormValue("password")
	data := fmt.Sprintf(`{"identity":"%s","password":"%s","passwordConfirm":"%s"}`, email, password, password)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/users/auth-with-password", strings.NewReader(data))
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
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
			return e.String(http.StatusInternalServerError, err.Error())
		}
		token_cookie := new(http.Cookie)
		token_cookie.Name = "explore_token"
		token_cookie.Value = authResponse.Token
		token_cookie.Secure = true
		token_cookie.HttpOnly = true
		token_cookie.Expires = time.Now().Add(1 * time.Hour)
		e.SetCookie(token_cookie)
		id_cookie := new(http.Cookie)
		id_cookie.Name = "id"
		id_cookie.Value = authResponse.Record.ID
		id_cookie.Secure = true
		id_cookie.HttpOnly = true
		id_cookie.Expires = time.Now().Add(1 * time.Hour)
		e.SetCookie(id_cookie)
		e.Response.Header().Set("HX-Redirect", "/")
	}
	loginError := "Login Failed, Please try again"
	return e.HTML(200, loginError)
}

func logoutHandler(e *core.RequestEvent) error {
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

func registerHandler(e *core.RequestEvent) error {
	email := e.Request.FormValue("email")
	password := e.Request.FormValue("password")
	data := fmt.Sprintf(`{"email":"%s","password":"%s","passwordConfirm":"%s"}`, email, password, password)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/users/records", strings.NewReader(data))
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		return e.Redirect(http.StatusMovedPermanently, "/")
	}

	return e.String(http.StatusBadRequest, string(body))
}

func validateToken(e *core.RequestEvent, app *pocketbase.PocketBase) error {
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
