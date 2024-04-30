package entity

import "github.com/Uikola/carCatalog/pkg/filter"

type ListCarsOptions struct {
	FilterOptions Options
	Limit         int
	Offset        int
}

type Options interface {
	IsToApply() bool
	AddField(name, operator, value, dType string) error
	Fields() []filter.Field
}
