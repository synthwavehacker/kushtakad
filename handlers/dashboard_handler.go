package handlers

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/kushtaka/kushtakad/events"
	"github.com/kushtaka/kushtakad/state"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Error(err)
	}

	var events []events.EventManager
	app.DB.All(&events, storm.Reverse())
	log.Info(events)
	app.View.Events = events
	app.View.Links.Dashboard = "active"
	app.Render.HTML(w, http.StatusOK, "admin/pages/dashboard", app.View)
	return
}
