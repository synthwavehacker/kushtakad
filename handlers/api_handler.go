package handlers

import (
	"net/http"
	"strings"

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
	if sensor.ApiKey != apiKey {
		app.Render.JSON(w, 404, err)
		return
	}

	svm := sensor.ServicesConfig(app.DB)

	app.Render.JSON(w, http.StatusOK, svm)
	return
}
