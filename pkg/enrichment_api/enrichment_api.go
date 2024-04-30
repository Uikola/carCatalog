package enrichment_api

import (
	"encoding/json"
	"fmt"
	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}

type GetCarByRegNumResp struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year,omitempty"`
	Owner  struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic,omitempty"`
	} `json:"owner"`
}

func (c Client) GetCarByRegNum(regNum string) (entity.Car, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/info?regNum=%s", os.Getenv("API_URL"), regNum), nil)
	if err != nil {
		return entity.Car{}, err
	}

	log.Info().Msgf("Making request for enrichment api")
	resp, err := c.client.Do(req)
	if err != nil {
		return entity.Car{}, err
	}
	defer resp.Body.Close()
	log.Info().Msgf("Get response from enrichment api. Status: %s, body: %v", resp.Status, resp.Body)

	var response GetCarByRegNumResp
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return entity.Car{}, err
	}

	car := entity.Car{
		RegNum: response.RegNum,
		Mark:   response.Mark,
		Model:  response.Model,
		Year:   response.Year,
		Owner: entity.People{
			Name:       response.Owner.Name,
			Surname:    response.Owner.Surname,
			Patronymic: response.Owner.Patronymic,
		},
	}

	return car, nil
}
