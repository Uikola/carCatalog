package entity

import (
	"github.com/Uikola/carCatalog/internal/errorz"
	"github.com/rs/zerolog/log"
	"regexp"
)

type UpdateCarRequest struct {
	RegNum string `json:"reg_num,omitempty"`
	Mark   string `json:"mark,omitempty"`
	Model  string `json:"model,omitempty"`
	Year   int    `json:"year,omitempty"`
	Owner  struct {
		Name       string `json:"name,omitempty"`
		Surname    string `json:"surname,omitempty"`
		Patronymic string `json:"patronymic,omitempty"`
	} `json:"owner,omitempty"`
}

func (uc UpdateCarRequest) Valid() error {
	if uc.RegNum != "" {
		re := regexp.MustCompile(`^[A-Z]\d{3}[A-Z]{2}\d{3}$`)
		if !re.MatchString(uc.RegNum) {
			return errorz.ErrInvalidRegNum
		}
	}
	return nil
}

type AddCarsRequest struct {
	RegNums []string `json:"regNums"`
}

func (ac AddCarsRequest) Valid() []error {
	var errs []error

	re := regexp.MustCompile(`^[A-Z]\d{3}[A-Z]{2}\d{3}$`)

	for _, regNum := range ac.RegNums {
		if !re.MatchString(regNum) {
			log.Error().Err(errorz.ErrInvalidRegNum).Msg("invalid registration number")
		}
	}
	return errs
}
