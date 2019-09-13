package handlers

import (
	"net/http"

	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetSmtp(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Error(err)
	}

	var smtp models.Smtp
	err = app.DB.One("ID", 1, &smtp)
	if err == nil {
		app.View.Forms.Smtp = &smtp

	}

	app.View.AddCrumb("SMTP", "#")
	app.Render.HTML(w, http.StatusOK, "admin/pages/smtp", app.View)
	return
}

func PostSmtp(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Error(err)
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
		http.Redirect(w, r, "/kushtaka/smtp", 302)
		return
	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/kushtaka/smtp", 302)
		return
	}

	err = tx.Save(smtp)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/kushtaka/smtp", 302)
		return
	}

	err = tx.Commit()
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/kushtaka/smtp", 302)
		return
	}

	app.Success("SMTP saved successfully.")
	http.Redirect(w, r, "/kushtaka/smtp", 302)
	return
}

func PutSmtp(w http.ResponseWriter, r *http.Request) {
	log.Error("PutSmtp()")
	return
}

func DeleteSmtp(w http.ResponseWriter, r *http.Request) {
	log.Error("DeleteSmtp()")
	return
}
