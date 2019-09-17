package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kushtaka/kushtakad/helpers"
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
	defer tx.Rollback()

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

func PostSendTestEmail(w http.ResponseWriter, r *http.Request) {
	log.Debug("PostSendTestEmail")
	w.Header().Set("Content-Type", "application/json")
	app, err := state.Restore(r)
	if err != nil {
		log.Fatalf("Unable to restore app : %s", err)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		app.Render.JSON(w, 200, err)
		return
	}
	defer r.Body.Close()

	smtp := &models.Smtp{}
	err = json.Unmarshal(b, smtp)
	if err != nil {
		log.Error(err)
		app.Render.JSON(w, 200, err)
		return
	}

	m := helpers.NewMailer(app.DB, app.Box)
	m.Smtp = smtp
	err = m.SendSensorEvent("1", "furl", "1", "new", "jfolkins@gmail.com", time.Now())
	if err != nil {
		log.Errorf("Failed to send email %s", err)
		app.Render.JSON(w, 200, NewResponse("failed", "Failed to send email", err))
		return
	}
	resp := &Response{}
	resp.Status = "success"
	resp.Message = "Email Sent"
	w.Write(resp.JSON())
	return
}
