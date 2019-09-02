package models

import (
	"log"

	"github.com/asdine/storm"
)

type State struct {
	IsAuthd      bool
	AdminIsSetup bool
	SmtpIsSetup  bool
	db           *storm.DB
}

func NewState(u *User, db *storm.DB) *State {
	st := &State{db: db}
	st.AdminIsSetup = st.isSetup()
	st.IsAuthd = st.isAuthd(u)
	return st
}

func (st *State) isSetup() bool {
	var user User
	err := st.db.One("ID", 1, &user)
	if err != nil {
		return false
	}

	if user.ID == 1 {
		return true
	}

	return false
}

func (st *State) isAuthd(u *User) bool {
	log.Println(u.ID)
	if u.ID > 0 {
		return true
	}

	return false
}
