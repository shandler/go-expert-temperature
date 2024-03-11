package domain

import (
	"net/http"
)

type Application interface {
	Run() error
}

type Action interface {
	Path() string
	Handle(w http.ResponseWriter, r *http.Request)
}
