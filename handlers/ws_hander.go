package handlers


import (
	"net/http"
	"log"
	"github.com/kushtaka/kushtakad/state"
)

func Ws(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
	}
	app.ServerHub.ServeWs(w, r)
}