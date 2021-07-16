package web

import (
	"encoding/gob"

	"github.com/google/uuid"
)

func init() {
	gob.Register(CreateToolForm{})
	gob.Register(CreateToolEntryForm{})
	gob.Register(RegisterForm{})
	gob.Register(LoginForm{})
	gob.Register(FormErrors{})
}

type FormErrors map[string]string

type CreatePostForm struct {
	Title   string
	Content string

	Errors FormErrors
}

type CreateToolForm struct {
	Name     string
	Model    string
	Category string

	Errors FormErrors
}

func (f *CreateToolForm) Validate() bool {
	f.Errors = FormErrors{}

	if f.Name == "" {
		f.Errors["Name"] = "Please enter a name."
	}

	return len(f.Errors) == 0
}

type CreateToolEntryForm struct {
	ToolID    uuid.UUID
	Condition string

	Errors FormErrors
}

func (f *CreateToolEntryForm) Validate() bool {
	f.Errors = FormErrors{}

	_, err := uuid.Parse(string(f.ToolID.String()))
	if err != nil {
		f.Errors["ToolID"] = "Please select a tool"
	}

	if f.Condition == "" {
		f.Errors["Condition"] = "Please select a condition."
	}

	return len(f.Errors) == 0
}

type RegisterForm struct {
	Username      string
	Password      string
	UsernameTaken bool

	Errors FormErrors
}

func (f *RegisterForm) Validate() bool {
	f.Errors = FormErrors{}

	if f.Username == "" {
		f.Errors["Username"] = "Please enter a username."
	} else if f.UsernameTaken {
		f.Errors["Username"] = "This username is already taken."
	}

	if f.Password == "" {
		f.Errors["Password"] = "Please enter a password."
	} else if len(f.Password) < 8 {
		f.Errors["Password"] = "Your password must be at least 8 characters long."
	}

	return len(f.Errors) == 0
}

type LoginForm struct {
	Username             string
	Password             string
	IncorrectCredentials bool

	Errors FormErrors
}

func (f *LoginForm) Validate() bool {
	f.Errors = FormErrors{}

	if f.Username == "" {
		f.Errors["Username"] = "Please enter a username."
	} else if f.IncorrectCredentials {
		f.Errors["Username"] = "Username or password is incorrect."
	}

	if f.Password == "" {
		f.Errors["Password"] = "Please enter a password."
	}

	return len(f.Errors) == 0
}
