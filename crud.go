package main

import (
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func saveLocation(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	//validate auth token before anything
	err := validateToken(e, app)
	if err != nil {
		return e.String(http.StatusOK, "Looks like you do not have a valid token, please login to save a location.")
	}

	//get location and id from request.
	location := e.Request.FormValue("location")
	id_cookie, err := e.Request.Cookie("id")
	id := id_cookie.Value

	//specfiy collection
	collection, err := app.FindCollectionByNameOrId("favorite_locations")
	if err != nil {
		return err
	}

	//create record programtically instead of with http request
	record := core.NewRecord(collection)
	record.Set("location", location)
	record.Set("user_id", id)
	err = app.Save(record)
	if err != nil {
		return err
	}

	//return data. TODO return everything in collection with places view.
	return e.Redirect(303, "/my_places")
}

func deleteLocation(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	//validate auth token before anything
	err := validateToken(e, app)
	if err != nil {
		return e.String(http.StatusOK, "Looks like you do not have a valid token, please login to delete a location")
	}

	//get the location from the request
	location := e.Request.URL.Query().Get("location")
	if location == "" {
		return e.String(http.StatusBadRequest, "Location not provided")
	}

	//get user's id from the request
	id_cookie, err := e.Request.Cookie("id")
	id := id_cookie.Value

	record, err := app.FindFirstRecordByFilter(
		"favorite_locations",
		"location = {:location} && user_id = {:user_id}",
		dbx.Params{"location": location, "user_id": id},
	)

	err = app.Delete(record)
	if err != nil {
		return err
	}

	//return data. TODO return everything in collection with places view.
	return e.Redirect(303, "/my_places")
}
