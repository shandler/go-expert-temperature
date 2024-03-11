package search_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shandler/go-expert-temperature/internal/config"
	"github.com/shandler/go-expert-temperature/internal/core/search"
	"github.com/shandler/go-expert-temperature/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type zipCodeHttpStub struct {
	mock.Mock
}

func (z *zipCodeHttpStub) Do(request *http.Request) (*http.Response, error) {
	args := z.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestZipCodeHttp(t *testing.T) {
	t.Run("should return error when zipCode is invalid", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		// Act
		req := dto.SearchRequest{ZipCode: "123456-789"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusUnprocessableEntity, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "invalid zipCode"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error when zipCode is not found", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		clientStub.On("Do", mock.Anything).Return(&http.Response{}, errors.New("can not found zipCode"))

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error when body is nil", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		clientStub.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{-}")),
		}, nil)

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error when body not localidade", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		requestMatch, _ := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/12345-678/json/", nil)

		clientStub.On("Do", requestMatch).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"zipCode": "12345678"}`)),
		}, nil)

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return the necessary information", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		requestMatch, _ := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/12345-678/json/", nil)
		clientStub.On("Do", requestMatch).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"localidade": "Sao Paulo"}`)),
		}, nil)

		requestWeather, _ := http.NewRequest(http.MethodGet, "http://api.weatherapi.com/v1/current.json?key=&q=Sao+Paulo", nil)
		clientStub.On("Do", requestWeather).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"current": {"temp_c": 10, "temp_f": 50}}`)),
		}, nil)

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		expectedBody := &search.WeatherResponse{
			TempC: 10,
			TempF: 50,
			TempK: 283,
		}

		assert.Equal(t, http.StatusOK, response.Status)
		assert.EqualValues(t, expectedBody, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error parse json in weather api", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		requestMatch, _ := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/12345-678/json/", nil)
		clientStub.On("Do", requestMatch).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"localidade": "Sao Paulo"}`)),
		}, nil)

		requestWeather, _ := http.NewRequest(http.MethodGet, "http://api.weatherapi.com/v1/current.json?key=&q=Sao+Paulo", nil)
		clientStub.On("Do", requestWeather).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"}}`)),
		}, nil)

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode in weather api"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error when statusCode can not found in weather api", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		requestMatch, _ := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/12345-678/json/", nil)
		clientStub.On("Do", requestMatch).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"localidade": "Sao Paulo"}`)),
		}, nil)

		requestWeather, _ := http.NewRequest(http.MethodGet, "http://api.weatherapi.com/v1/current.json?key=&q=Sao+Paulo", nil)
		clientStub.On("Do", requestWeather).Return(&http.Response{}, errors.New("can not found zipCode in weather api"))

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode in weather api"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error when current can not found in weather api", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		requestMatch, _ := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/12345-678/json/", nil)
		clientStub.On("Do", requestMatch).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"localidade": "Sao Paulo"}`)),
		}, nil)

		requestWeather, _ := http.NewRequest(http.MethodGet, "http://api.weatherapi.com/v1/current.json?key=&q=Sao+Paulo", nil)
		clientStub.On("Do", requestWeather).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		}, nil)

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode in weather api"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error when current.temp_c can not found in weather api", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		requestMatch, _ := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/12345-678/json/", nil)
		clientStub.On("Do", requestMatch).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"localidade": "Sao Paulo"}`)),
		}, nil)

		requestWeather, _ := http.NewRequest(http.MethodGet, "http://api.weatherapi.com/v1/current.json?key=&q=Sao+Paulo", nil)
		clientStub.On("Do", requestWeather).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"current": {}}`)),
		}, nil)

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode in weather api"}, response.Body)

		clientStub.AssertExpectations(t)
	})

	t.Run("should return error when current.temp_f can not found in weather api", func(t *testing.T) {
		// Arrange
		clientStub := &zipCodeHttpStub{}
		config := config.New()

		zcr := search.New(clientStub, config)

		requestMatch, _ := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/12345-678/json/", nil)
		clientStub.On("Do", requestMatch).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"localidade": "Sao Paulo"}`)),
		}, nil)

		requestWeather, _ := http.NewRequest(http.MethodGet, "http://api.weatherapi.com/v1/current.json?key=&q=Sao+Paulo", nil)
		clientStub.On("Do", requestWeather).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"current": { "temp_c": 10 }}`)),
		}, nil)

		// Act
		req := dto.SearchRequest{ZipCode: "12345678"}
		response := zcr.Search(req)

		// Assert
		assert.Equal(t, http.StatusNotFound, response.Status)
		assert.EqualValues(t, struct{ Message string }{Message: "can not found zipCode in weather api"}, response.Body)

		clientStub.AssertExpectations(t)
	})
}
