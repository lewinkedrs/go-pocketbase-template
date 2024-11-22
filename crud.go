package main

import (
	"fmt"
	"github.com/labstack/echo/v5"
	"io"
	"net/http"
	"strings"
)

func saveLocation(c echo.Context) error {
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
}
