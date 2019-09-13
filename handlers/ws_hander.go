package handlers


import (
	"net/http"
	"github.com/kushtaka/kushtakad/state"
)

func Ws(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Error(err)
	}
	app.ServerHub.ServeWs(w, r)
}