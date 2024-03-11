package actions

import (
	"net/http"
)

type actionDefault struct {
	path string
}

func NewDefault(path string) *actionDefault {
	return &actionDefault{
		path: path,
	}
}

func (a actionDefault) Path() string {
	return a.path
}

func (a actionDefault) Handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("access: /zip-code?zipCode=07987110"))
}
