package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/shandler/go-expert-temperature/internal/config"
	"github.com/shandler/go-expert-temperature/internal/domain"
	"github.com/shandler/go-expert-temperature/internal/dto"
)

type zipCodeHttp struct {
	client domain.HTTPClient
	config *config.Config
}

type WeatherResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func New(client domain.HTTPClient, config *config.Config) *zipCodeHttp {
	return &zipCodeHttp{
		client: client,
		config: config,
	}
}

func (z *zipCodeHttp) Search(request dto.SearchRequest) dto.SearchResponse {

	regex := regexp.MustCompile("^[0-9]{8}$")
	if !regex.MatchString(request.ZipCode) {
		return z.mountError(http.StatusUnprocessableEntity, "invalid zipCode")
	}

	locale, err := z.FindZipCode(request.ZipCode)
	if err != nil {
		return z.mountError(http.StatusNotFound, err.Error())
	}

	response, err := z.FindWeather(locale)
	if err != nil {
		return z.mountError(http.StatusNotFound, err.Error())
	}

	return dto.SearchResponse{
		Status: http.StatusOK,
		Body:   response,
	}
}

func (z *zipCodeHttp) mountError(status int, err string) dto.SearchResponse {
	return dto.SearchResponse{
		Status: status,
		Body: struct {
			Message string `json:"message"`
		}{Message: err},
	}
}

func (z *zipCodeHttp) FindZipCode(zipCode string) (string, error) {
	urlZipCode := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipCode[0:5]+"-"+zipCode[5:])
	log.Println(urlZipCode)

	requestViaCep, _ := http.NewRequest(http.MethodGet, urlZipCode, nil)

	response, err := z.client.Do(requestViaCep)
	if err != nil {
		log.Println(err)
		return "", errors.New("can not found zipCode")
	}

	var body map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		log.Println(err)
		return "", errors.New("can not found zipCode")
	}

	if _, ok := body["localidade"].(string); !ok {
		log.Println("localidade not found")
		return "", errors.New("can not found zipCode")
	}

	return body["localidade"].(string), nil
}

func (z *zipCodeHttp) FindWeather(locale string) (*WeatherResponse, error) {
	urlWeather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", z.config.WeatherKey, url.QueryEscape(locale))
	requestWeather, _ := http.NewRequest(http.MethodGet, urlWeather, nil)

	response, err := z.client.Do(requestWeather)
	if err != nil {
		return nil, errors.New("can not found zipCode in weather api")
	}

	var body map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, errors.New("can not found zipCode in weather api")
	}

	current, ok := body["current"]
	if !ok {
		return nil, errors.New("can not found zipCode in weather api")
	}

	tempC, ok := current.(map[string]interface{})["temp_c"].(float64)
	if !ok {
		return nil, errors.New("can not found zipCode in weather api")
	}

	tempF, ok := current.(map[string]interface{})["temp_f"].(float64)
	if !ok {
		return nil, errors.New("can not found zipCode in weather api")
	}

	tempK := tempC + 273

	return &WeatherResponse{TempC: tempC, TempF: tempF, TempK: tempK}, nil
}
