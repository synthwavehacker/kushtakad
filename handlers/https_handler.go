package handlers

import (
	"net/http"

	"github.com/kushtaka/kushtakad/state"
)

func GetHttps(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		app.Render.JSON(w, 404, err)
		return
	}

	app.View.AddCrumb("HTTPS", "#")
	app.Render.HTML(w, http.StatusOK, "admin/pages/https", app.View)
	return
}

func PostHttps(w http.ResponseWriter, r *http.Request) {
	log.Debug("PostHttps")
	url := "/kushtaka/https"
	app, err := state.Restore(r)
	if err != nil {
		app.Render.JSON(w, 404, err)
		return
	}

	app.Reboot <- true

	app.Success("Reboot started...")
	http.Redirect(w, r, url, 302)
	return
}
