package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetTokens(w http.ResponseWriter, r *http.Request) {
	redirUrl := "/kushtaka/dashboard"
	app, err := state.Restore(r)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/404", 404)
		return
	}

	var tokens []models.Token
	err = app.DB.All(&tokens)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	app.View.Tokens = tokens
	app.Render.HTML(w, http.StatusOK, "admin/pages/tokens", app.View)
	return
}

func PostTokens(w http.ResponseWriter, r *http.Request) {
	redirUrl := "/kushtaka/tokens/page/1/limit/100"
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	token := &models.Token{Name: name}
	app.View.Forms.Token = token
	err = token.ValidateCreate()
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	tx.One("Name", name, token)
	if token.ID > 0 {
		tx.Rollback()
		app.Fail("Token using that name already exists.")
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	err = tx.Save(token)
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

	app.View.Forms = state.NewForms()
	app.Success(fmt.Sprintf("The token [%s] was created successfully.", token.Name))
	http.Redirect(w, r, redirUrl, 301)
	return
}
