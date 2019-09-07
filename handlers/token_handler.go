package handlers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
	"github.com/kushtaka/kushtakad/tokens/pdf"
)

func GetTestToken(w http.ResponseWriter, r *http.Request) {
	log.Println("test token")
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
		return
	}

	i, err := app.Box.Find("files/i.png")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(i)))
	http.ServeContent(w, r, "canary.docx", time.Now(), bytes.NewReader(i))
}

func CreateDocxToken(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
		return
	}

	token := &models.Token{}
	ts, err := token.BuildDocx(app.Box)
	if err != nil {
		log.Println(err)
		return
	}

	data, err := ioutil.ReadFile(ts)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=canary.docx")
	http.ServeContent(w, r, "canary.docx", time.Now(), bytes.NewReader(data))
}

func CreatePdfToken(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
		return
	}

	pdfc, err := pdf.NewPdfContext("http://localhost:3000", app.Box)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("PDFC-URL", pdfc.Url)

	w.Header().Set("Content-Disposition", "attachment; filename=canary.pdf")
	http.ServeContent(w, r, "canary.pdf", time.Now(), bytes.NewReader(pdfc.Buffer.Bytes()))
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	redirUrl := "/kushtaka/teams/page/1/limit/100"
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	token := &models.Token{}
	err = app.DB.One("ID", id, token)
	if err != nil {
		app.Fail("Token does not exist")
		http.Redirect(w, r, redirUrl, 302)
		return
	}

	app.View.Token = token
	app.View.Links.Tokens = "active"
	app.View.AddCrumb("Tokens", "/kushtaka/tokens/page/1/limit/100")
	app.View.AddCrumb(token.Name, "#")
	app.Render.HTML(w, http.StatusOK, "admin/pages/token", app.View)
	return
}

func PostToken(w http.ResponseWriter, r *http.Request) {
	/*
		redirUrl := "/kushtaka/teams/page/1/limit/100"

		app, err := state.Restore(r)
		if err != nil {
			log.Println(err)
		}

		email := r.FormValue("email")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			app.Fail("Unable to parse ID")
			http.Redirect(w, r, redirUrl, 302)
			return
		}

		token := &models.Token{}
		err = app.DB.One("ID", id, token)
		if err != nil {
			app.Fail("Token does not exist. " + err.Error())
			http.Redirect(w, r, redirUrl, 302)
			return
		}

		url := "/kushtaka/token/" + vars["id"]
		err = team.ValidateAddMember(email)
		app.View.Forms.TeamMember = team
		if err != nil {
			app.Fail(err.Error())
			http.Redirect(w, r, url, 302)
			return
		}

		tx, err := app.DB.Begin(true)
		if err != nil {
			app.Fail(err.Error())
			http.Redirect(w, r, url, 302)
			return
		}
		team.MemberToAdd = ""

		err = tx.Save(team)
		if err != nil {
			app.Fail(err.Error())
			http.Redirect(w, r, url, 302)
			return
		}

		err = tx.Commit()
		if err != nil {
			app.Fail(err.Error())
			http.Redirect(w, r, "/kushtaka/dashboard", 302)
			return
		}

		app.View.Forms = state.NewForms()
		app.Success("Member has been successfully added to the team.")
		http.Redirect(w, r, url, 302)
	*/
	return
}

func PutToken(w http.ResponseWriter, r *http.Request) {
	log.Println("PutToken()")
	return
}

func DeleteToken(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteToken()")
	return
}
