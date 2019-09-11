package server

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"

	"github.com/asdine/storm"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kushtaka/kushtakad/handlers"
	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
	"github.com/urfave/negroni"
)

const assetsFolder = "static"
const sessionName = "_kushtaka"

var (
	rtr      *mux.Router
	fss      *sessions.FilesystemStore
	db       *storm.DB
	box      *packr.Box
	settings models.Settings
	err      error
)

func Run() {

	gob.Register(&state.App{})
	box = packr.New(assetsFolder, "../static")

	err = state.SetupFileStructure(box)
	if err != nil {
		log.Fatal(err)
	}

	db, err = storm.Open(state.DbLocation())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// must setup the basic hashes and settings for application to function
	settings, err = models.InitSettings(db)
	if err != nil {
		log.Fatal(err)
	}

	fss = sessions.NewFilesystemStore(state.SessionLocation(), settings.SessionHash, settings.SessionBlock)

	// open
	rtr = mux.NewRouter()
	rtr.HandleFunc("/assets/{theme}/{dir}/{file}", handlers.Asset).Methods("GET")
	rtr.HandleFunc("/setup", handlers.GetSetup).Methods("GET")
	rtr.HandleFunc("/setup", handlers.PostSetup).Methods("POST")
	rtr.HandleFunc("/t/{t}/i.png", handlers.GetTestToken).Methods("GET")
	rtr.HandleFunc("/create-pdf-token", handlers.CreatePdfToken).Methods("GET")
	rtr.HandleFunc("/create-docx-token", handlers.CreateDocxToken).Methods("GET")

	rtr.HandleFunc("/", handlers.IndexCheckr).Methods("GET")
	rtr.NotFoundHandler = &NotFound{}

	// login has its own middleware chain
	login := mux.NewRouter().PathPrefix("/login").Subrouter().StrictSlash(false)
	login.Use(forceSetup)
	login.HandleFunc("", handlers.GetLogin).Methods("GET")
	login.HandleFunc("", handlers.PostLogin).Methods("POST")

	// mod has its own middleware chain
	// protected, can't process unless logged in and setup is complete
	kushtaka := mux.NewRouter().PathPrefix("/kushtaka").Subrouter().StrictSlash(true)
	kushtaka.Use(forceSetup)
	kushtaka.Use(isAuthenticated)
	kushtaka.HandleFunc("/dashboard", handlers.GetDashboard).Methods("GET")

	// sensors
	kushtaka.HandleFunc("/sensors/page/{pid}/limit/{oid}", handlers.GetSensors).Methods("GET")
	kushtaka.HandleFunc("/sensors", handlers.PostSensors).Methods("POST")
	// sensor
	kushtaka.HandleFunc("/sensor/{id}", handlers.GetSensor).Methods("GET")
	kushtaka.HandleFunc("/sensor", handlers.PostSensor).Methods("POST")

	// service
	kushtaka.HandleFunc("/service/{sensor_id}/type/{type}", handlers.PostService).Methods("POST")
	kushtaka.HandleFunc("/service", handlers.DeleteService).Methods("DELETE")

	// tokens
	kushtaka.HandleFunc("/tokens/page/{pid}/limit/{oid}", handlers.GetTokens).Methods("GET")
	kushtaka.HandleFunc("/tokens", handlers.PostTokens).Methods("POST")
	// token
	kushtaka.HandleFunc("/token/{id}", handlers.GetToken).Methods("GET")
	kushtaka.HandleFunc("/token", handlers.PostToken).Methods("POST")
	kushtaka.HandleFunc("/token", handlers.PutToken).Methods("PUT")
	kushtaka.HandleFunc("/token", handlers.DeleteToken).Methods("DELETE")

	// smtp
	kushtaka.HandleFunc("/smtp", handlers.GetSmtp).Methods("GET")
	kushtaka.HandleFunc("/smtp", handlers.PostSmtp).Methods("POST")

	// users
	kushtaka.HandleFunc("/users/page/{pid}/limit/{oid}", handlers.GetUsers).Methods("GET")
	kushtaka.HandleFunc("/users", handlers.PostUsers).Methods("POST")

	// user
	kushtaka.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	kushtaka.HandleFunc("/user/{id}", handlers.PostUser).Methods("POST")
	kushtaka.HandleFunc("/user/{id}", handlers.PutUser).Methods("PUT")
	kushtaka.HandleFunc("/user/{id}", handlers.DeleteUser).Methods("DELETE")

	// teams
	kushtaka.HandleFunc("/teams/page/{pid}/limit/{oid}", handlers.GetTeams).Methods("GET")
	kushtaka.HandleFunc("/teams", handlers.PostTeams).Methods("POST")
	// team
	kushtaka.HandleFunc("/team/{id}", handlers.GetTeam).Methods("GET")
	kushtaka.HandleFunc("/team/{id}", handlers.PostTeam).Methods("POST")
	kushtaka.HandleFunc("/team/{id}", handlers.PutTeam).Methods("PUT")
	kushtaka.HandleFunc("/team/{id}", handlers.DeleteTeam).Methods("DELETE")

	// settings
	//kushtaka.HandleFunc("/settings", handlers.ReadSettings).Methods("GET")
	//kushtaka.HandleFunc("/settings", handlers.ReadSettings).Methods("POST")

	// wire up sub routers
	rtr.PathPrefix("/login").Handler(negroni.New(
		negroni.Wrap(login),
	))

	rtr.PathPrefix("/kushtaka").Handler(negroni.New(
		negroni.Wrap(kushtaka),
	))

	rtr.HandleFunc("/ws", handlers.Ws)

	// setup router
	n := negroni.New()
	n.Use(negroni.HandlerFunc(before))
	n.UseHandler(rtr)
	n.Use(negroni.HandlerFunc(after))

	log.Fatal(http.ListenAndServe(settings.Host, n))
}

// forceSetup is a middleware function that makes sure
// a admin user is created before allowing the user to
// proceed with using the application
func forceSetup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app := r.Context().Value(state.AppStateKey).(*state.App)
		var user models.User
		err := db.One("ID", 1, &user)
		if err != nil && r.URL.Path != "/setup" {
			app.Fail("You must create an admin user before proceeding.")
			http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app := r.Context().Value(state.AppStateKey).(*state.App)
		if app.User.ID < 1 {
			app.Fail("You must login before proceeding.")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func before(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// setup session and if it errors, create a new session
	sess, err := fss.Get(r, sessionName)
	if err != nil {
		fss.New(r, sessionName)
		sess, err = fss.Get(r, sessionName)
	}
	sess.Options.HttpOnly = true

	app, err := state.NewApp(w, r, db, sess, fss, box)
	if err != nil {
		log.Fatal(err)
	}
	app.RestoreFlash()
	app.RestoreUser()
	app.RestoreForm()
	app.RestoreState()
	app.RestoreURI()

	ctx := context.WithValue(r.Context(), state.AppStateKey, app)
	next(w, r.WithContext(ctx))
}

func after(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	app := r.Context().Value(state.AppStateKey).(*state.App)

	// because we build the view upon each request
	// we clear it here to keep consistency and state
	//

	userState, err := json.Marshal(app.User)
	if err != nil {
		log.Fatal(err)
	}

	formState, err := json.Marshal(app.View.Forms)
	if err != nil {
		log.Fatal(err)
	}

	app.Session.Values[state.UserStateKey] = userState
	app.Session.Values[state.FormStateKey] = formState
	err = app.Session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}

	app.View.Clear()

	next(w, r)
}

//
// NOT FOUND
//
type NotFound struct{}

func (nf *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 Not Found"))
	return
}
