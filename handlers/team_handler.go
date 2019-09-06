package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetTeam(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	team := &models.Team{}
	err = app.DB.One("ID", id, team)
	if err != nil {
		app.Fail("Team does not exist")
		http.Redirect(w, r, "/kushtaka/teams/page/1/limit/100", 302)
		return
	}

	app.View.Team = team
	app.Render.HTML(w, http.StatusOK, "admin/pages/team", app.View)
	return
}

func PostTeam(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
	}

	email := r.FormValue("email")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.Fail("Unable to parse ID")
		http.Redirect(w, r, "/kushtaka/teams/page/1/limit/100", 302)
		return
	}

	team := &models.Team{}
	err = app.DB.One("ID", id, team)
	if err != nil {
		app.Fail("Team does not exist. " + err.Error())
		http.Redirect(w, r, "/kushtaka/teams/page/1/limit/100", 302)
		return
	}

	url := "/kushtaka/team/" + vars["id"]
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
	return
}

func PutTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("PutTeam()")
	return
}

func DeleteTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteTeam()")
	return
}
