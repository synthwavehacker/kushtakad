package handlers

import (
	"log"
	"net/http"

	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUser()")
	return
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	log.Println("PostUser()")
	return
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	log.Println("PutUser()")
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteUser()")
	return
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	redirUrl := "/kushtaka/dashboard"
	app, err := state.Restore(r)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/404", 404)
		return
	}
	app.View.Links.Sensors = "active"

	var users []models.User
	err = app.DB.All(&users)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	app.View.Users = users
	app.Render.HTML(w, http.StatusOK, "admin/pages/users", app.View)
	return
}

func PostUsers(w http.ResponseWriter, r *http.Request) {
	redir := "/kushtaka/users/page/1/limit/100"
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
	}

	user := &models.User{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirm"),
	}

	err = user.ValidateCreateUser()
	app.View.Forms.User = user
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 301)
		return
	}

	user.HashPassword()

	tx, err := app.DB.Begin(true)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 301)
		return
	}

	err = tx.Save(user)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 301)
		return
	}

	err = tx.Commit()
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 301)
		return
	}

	app.View.Forms = state.NewForms()
	app.Success("User created successfully")
	http.Redirect(w, r, redir, 301)
	return
}
