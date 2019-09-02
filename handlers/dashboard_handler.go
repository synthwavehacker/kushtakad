package handlers

import (
	"log"
	"net/http"

	"github.com/kushtaka/kushtakad/state"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
	}

	app.View.Links.Dashboard = "active"
	app.Render.HTML(w, http.StatusOK, "admin/pages/endpoints", app.View)
	return
}
