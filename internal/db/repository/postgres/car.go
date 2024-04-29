package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Uikola/carCatalog/internal/entity"
	"github.com/Uikola/carCatalog/internal/errorz"
	"github.com/Uikola/carCatalog/pkg/filter"
	"github.com/jmoiron/sqlx"
)

type CarRepository struct {
	db *sqlx.DB
}

func NewCarRepository(db *sqlx.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (r CarRepository) ListCars(ctx context.Context, options entity.ListCarsOptions) ([]entity.Car, error) {
	const op = "CarRepository.ListCars"

	query := `
	SELECT id, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic
	FROM cars `
	var values []interface{}

	if options.FilterOptions.IsToApply() {
		query, values = addFilters(query, options.FilterOptions)
	}
	query = addPagination(query, len(options.FilterOptions.Fields()))
	values = append(values, options.Limit, options.Offset)

	stmt, err := r.db.Preparex(query)
	if err != nil {
		return []entity.Car{}, fmt.Errorf("%s:%w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, values...)
	if err != nil {
		return []entity.Car{}, fmt.Errorf("%s:%w", op, err)
	}

	var regNum, mark, model, ownerName, ownerSurname, ownerPatronymic string
	var id uint
	var year int
	var cars []entity.Car
	for rows.Next() {
		if err = rows.Scan(&id, &regNum, &mark, &model, &year, &ownerName, &ownerSurname, &ownerPatronymic); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		car := entity.Car{
			ID:     id,
			RegNum: regNum,
			Mark:   mark,
			Model:  model,
			Year:   year,
			Owner: entity.People{
				Name:       ownerName,
				Surname:    ownerSurname,
				Patronymic: ownerPatronymic,
			},
		}
		cars = append(cars, car)
	}

	if len(cars) == 0 {
		return []entity.Car{}, errorz.ErrNoCars
	}

	return cars, nil
}

func addFilters(query string, filterOptions entity.Options) (string, []interface{}) {
	numCondMap := map[string]string{
		filter.OperatorEq:               "=",
		filter.OperatorNotEq:            "!=",
		filter.OperatorGreaterThan:      ">",
		filter.OperatorGreaterThanEqual: ">=",
		filter.OperatorLowerThan:        "<",
		filter.OperatorLowerThanEqual:   "<=",
	}
	var filterValues []interface{}

	filterFields := filterOptions.Fields()
	query += "WHERE "
	for i, filterField := range filterFields {
		if filterField.Type == filter.DataTypeStr {
			query += fmt.Sprintf("%s %s $%d AND ", filterField.Name, filterField.Operator, i+1)
		} else if filterField.Type == filter.DataTypeInt {
			query += fmt.Sprintf("%s %s $%d AND ", filterField.Name, numCondMap[filterField.Operator], i+1)
		}
		filterValues = append(filterValues, filterField.Value)
	}
	query = query[:len(query)-5]
	return query, filterValues
}

func addPagination(query string, ind int) string {
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", ind+1, ind+2)
	return query
}

func (r CarRepository) DeleteCar(ctx context.Context, carID uint) error {
	const op = "CarRepository.DeleteCar"

	stmt, err := r.db.Preparex("DELETE FROM cars WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = stmt.ExecContext(ctx, carID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (r CarRepository) UpdateCar(ctx context.Context, carID uint, carMap map[string]interface{}) error {
	const op = "CarRepository.UpdateCar"

	var updateQuery string
	var args []interface{}
	for key, value := range carMap {
		if value.(string) == "" {
			continue
		}
		updateQuery += fmt.Sprintf("%s=?, ", key)
		args = append(args, value)
	}

	updateQuery = updateQuery[:len(updateQuery)-2] + " WHERE id = ?"
	args = append(args, carID)

	stmt, err := r.db.Preparex("UPDATE cars SET " + updateQuery)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (r CarRepository) Exists(ctx context.Context, carID uint) (bool, error) {
	const op = "CarRepository.Exists"

	stmt, err := r.db.Preparex("SELECT reg_num FROM cars WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, carID)

	var regNum string
	err = row.Scan(&regNum)
	switch {
	case errors.Is(sql.ErrNoRows, err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return true, nil

}

func (r CarRepository) AddCar(ctx context.Context, car entity.Car) (entity.Car, error) {
	const op = "CarRepository.AddCar"

	stmt, err := r.db.Preparex("INSERT INTO cars(reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic")
	if err != nil {
		return entity.Car{}, fmt.Errorf("%s:%w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic)

	var id uint
	var regNum, mark, model, ownerName, ownerSurname, ownerPatronymic string
	var year int
	err = row.Scan(&id, &regNum, &mark, &model, &year, &ownerName, &ownerSurname, &ownerPatronymic)
	if err != nil {
		return entity.Car{}, fmt.Errorf("%s:%w", op, err)
	}

	return entity.Car{
		ID:     id,
		RegNum: regNum,
		Mark:   mark,
		Model:  model,
		Year:   year,
		Owner: entity.People{
			Name:       ownerName,
			Surname:    ownerSurname,
			Patronymic: ownerPatronymic,
		},
	}, nil
}
