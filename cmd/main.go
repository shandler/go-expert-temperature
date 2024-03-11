package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shandler/go-expert-temperature/internal/config"
	"github.com/shandler/go-expert-temperature/internal/core/search"
	"github.com/shandler/go-expert-temperature/internal/domain"
	"github.com/shandler/go-expert-temperature/internal/infra/actions"
	"github.com/shandler/go-expert-temperature/internal/infra/server"
)

func main() {
	// config
	config := config.New()

	// services
	var clientHttp domain.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	var zipCode domain.ZipCode = search.New(clientHttp, config)

	// actions
	var actionDefault domain.Action = actions.NewDefault("/")
	var actionZipCode domain.Action = actions.NewZipCode("/zip-code", zipCode)

	// server
	var app domain.Application = server.New(":8080", []domain.Action{actionZipCode, actionDefault})

	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
