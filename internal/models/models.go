package models

import (
	"errors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/utils"
)

type User struct {
	Id        string `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   int    `json:"balance"`
	Spent     int    `json:"spent"`
}

func (u *User) Validate() []error {
	errorSlice := []error{}
	errorSlice = append(errorSlice, u.validatePassword()...)
	errorSlice = append(errorSlice, u.validateLogin()...)
	return errorSlice
}

func (u *User) validatePassword() []error {
	errorSlice := []error{}
	if len(u.Password) < 6 {
		errorSlice = append(errorSlice, errors.New("password must be at least 6 characters"))
	}
	if utils.IsNumeric(u.Password) {
		errorSlice = append(errorSlice, errors.New("password entirely numeric"))
	}
	return errorSlice
}

func (u *User) validateLogin() []error {
	errorSlice := []error{}
	if len(u.Login) < 6 {
		errorSlice = append(errorSlice, errors.New("login must be at least 6 characters"))
	}
	return errorSlice
}
