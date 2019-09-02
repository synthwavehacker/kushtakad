package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetTeams(w http.ResponseWriter, r *http.Request) {
	redirUrl := "/kushtaka/dashboard"
	app, err := state.Restore(r)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/404", 404)
		return
	}

	var teams []models.Team
	err = app.DB.All(&teams)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	app.View.Teams = teams
	app.Render.HTML(w, http.StatusOK, "admin/pages/teams", app.View)
	return
}

func PostTeams(w http.ResponseWriter, r *http.Request) {
	redirUrl := "/kushtaka/teams/page/1/limit/100"
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	team := &models.Team{Name: name}
	app.View.Forms.Team = team
	err = team.ValidateCreate()
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	tx.One("Name", name, team)
	if team.ID > 0 {
		tx.Rollback()
		app.Fail("Team using that name already exists.")
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	err = tx.Save(team)
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	app.View.Forms.Team = models.NewTeam()
	app.Success(fmt.Sprintf("The team [%s] was created successfully.", name))
	http.Redirect(w, r, redirUrl, 301)
	return
}
