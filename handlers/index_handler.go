package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/kushtaka/kushtakad/state"
)

func IndexCheckr(w http.ResponseWriter, r *http.Request) {
	log.Debug("Validating schema and rebuilding indexes")
	http.Redirect(w, r, "/login", 302)
	return
}

func Asset(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Error(err)
	}

	v := mux.Vars(r)
	fp := filepath.Join(v["theme"], "assets", v["dir"], v["file"])
	b, err := app.Box.Find(fp)
	if err != nil {
		http.Error(w, fmt.Sprintf("404 Not Found: %s", v["file"]), 404)
		return
	}

	contentType := mime.TypeByExtension(path.Ext(fp))
	if contentType == "" {
		contentType = http.DetectContentType(b)
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(b)
	return
}
