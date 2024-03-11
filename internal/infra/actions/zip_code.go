package actions

import (
	"encoding/json"
	"net/http"

	"github.com/shandler/go-expert-temperature/internal/domain"
	"github.com/shandler/go-expert-temperature/internal/dto"
)

type actionZipCode struct {
	path    string
	zipCode domain.ZipCode
}

func NewZipCode(path string, zipCode domain.ZipCode) *actionZipCode {
	return &actionZipCode{
		path:    path,
		zipCode: zipCode,
	}
}

func (a actionZipCode) Path() string {
	return a.path
}

func (a actionZipCode) Handle(w http.ResponseWriter, r *http.Request) {

	request := dto.SearchRequest{ZipCode: r.URL.Query().Get("zipCode")}
	response := a.zipCode.Search(request)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response.Body)
}
