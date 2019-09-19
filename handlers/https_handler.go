package handlers

import (
	"net"
	"net/http"
	"sync"

	"github.com/kushtaka/kushtakad/state"
	"github.com/mholt/certmagic"
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
	log.Debug("POst Https")
	app, err := state.Restore(r)
	if err != nil {
		resp := NewResponse("failed", "failed to restore", err)
		app.Render.JSON(w, 200, resp)
		return
	}

	conn, err := net.Listen("tcp", "localhost:80")
	if err != nil {
		resp := NewResponse("failed", "Unable to bind to port 80", err)
		app.Render.JSON(w, 200, resp)
		return
	}
	conn.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		magic := certmagic.NewDefault()
		magic.CA = certmagic.LetsEncryptStagingCA
		magic.Email = "jfolkins@gmail.com"
		magic.Agreed = true
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			log.Debug("Lookit my cool website over HTTPS!")
			wg.Done()
		})
		err = http.ListenAndServe(":80", magic.HTTPChallengeHandler(mux))
		if err != nil {
			log.Debug(err)
		}
	}()

	wg.Wait()

	//app.Reboot <- true

	resp := NewResponse("success", "Port 80 is open", nil)
	app.Render.JSON(w, 200, resp)
	return
}
