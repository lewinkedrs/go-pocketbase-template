package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

func loginHandler(c echo.Context) error {
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
}

func logoutHandler(c echo.Context) error {
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
}

func registerHandler(c echo.Context) error {
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
}
