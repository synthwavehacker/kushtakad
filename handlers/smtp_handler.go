package handlers

import (
	"log"
	"net/http"

	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetSmtp(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
	}

	var smtp models.Smtp
	err = app.DB.One("ID", 1, &smtp)
	if err == nil {
		app.View.Forms.Smtp = &smtp

	}

	app.Render.HTML(w, http.StatusOK, "admin/pages/smtp", app.View)
	return
}

func PostSmtp(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Println(err)
	}

	smtp := &models.Smtp{
		ID:       1,
		Sender:   r.FormValue("sender"),
		Email:    r.FormValue("email"),
		Host:     r.FormValue("host"),
		Port:     r.FormValue("port"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	err = smtp.ValidateSmtp()
	app.View.Forms.Smtp = smtp
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/kushtaka/smtp", 301)
		return
	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/kushtaka/smtp", 301)
		return
	}

	err = tx.Save(smtp)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/kushtaka/smtp", 301)
		return
	}

	err = tx.Commit()
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/kushtaka/smtp", 301)
		return
	}

	app.Success("Smtp saved successfully.")
	http.Redirect(w, r, "/kushtaka/smtp", 301)
	return
}

func PutSmtp(w http.ResponseWriter, r *http.Request) {
	log.Println("PutSmtp()")
	return
}

func DeleteSmtp(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteSmtp()")
	return
}
