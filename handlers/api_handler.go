package handlers

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/kushtaka/kushtakad/events"
	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetConfig(w http.ResponseWriter, r *http.Request) {
	var sensor models.Sensor
	var apiKey string
	app, err := state.Restore(r)
	if err != nil {
		app.Render.JSON(w, 404, err)
		return
	}

	token, ok := r.Header["Authorization"]
	if ok && len(token) >= 1 {
		apiKey = token[0]
		apiKey = strings.TrimPrefix(apiKey, "Bearer ")
	}

	app.DB.One("ApiKey", apiKey, &sensor)
	// TODO: add constant time compare
	// update: not needed, handled in middleware
	if sensor.ApiKey != apiKey {
		app.Render.JSON(w, 404, err)
		return
	}

	svm := sensor.ServicesConfig(app.DB)

	app.Render.JSON(w, http.StatusOK, svm)
	return
}

func PostEvent(w http.ResponseWriter, r *http.Request) {
	var sensor models.Sensor
	var apiKey string
	app, err := state.Restore(r)
	if err != nil {
		app.Render.JSON(w, 404, err)
		return
	}

	token, ok := r.Header["Authorization"]
	if ok && len(token) >= 1 {
		apiKey = token[0]
		apiKey = strings.TrimPrefix(apiKey, "Bearer ")
	}

	// TODO: add constant time compare
	// update: not needed, handled in middleware
	app.DB.One("ApiKey", apiKey, &sensor)
	if sensor.ApiKey != apiKey {
		app.Render.JSON(w, 404, err)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		app.Render.JSON(w, 404, err)
		return
	}
	defer r.Body.Close()

	var em events.EventManager
	err = json.Unmarshal(b, &em)
	if err != nil {
		log.Error(err)
		app.Render.JSON(w, 404, err)
		return
	}

	err = app.DB.Save(&em)
	if err != nil {
		log.Error(err)
		app.Render.JSON(w, 404, err)
		return
	}

	app.Render.JSON(w, http.StatusOK, "success")
	return
}
