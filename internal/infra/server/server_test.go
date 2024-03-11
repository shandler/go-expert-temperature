package server_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/shandler/go-expert-temperature/internal/domain"
	"github.com/shandler/go-expert-temperature/internal/infra/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ActionMock struct {
	mock.Mock
}

func (a *ActionMock) Path() string {
	args := a.Called()
	return args.String(0)
}

func (a *ActionMock) Handle(w http.ResponseWriter, r *http.Request) {
	a.Called(w, r)
}

var actionStub *ActionMock = &ActionMock{}

func TestServerNew(t *testing.T) {

	server := server.New(":8080", []domain.Action{actionStub})
	assert.NotNil(t, server)

	// test run server
	actionStub.On("Path").Return("/test")
	actionStub.On("Handle", mock.Anything, mock.Anything).Return()

	go server.Run()
	time.Sleep(100 * time.Millisecond)

	res, err := http.Get("http://127.0.0.1:8080/test")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}
