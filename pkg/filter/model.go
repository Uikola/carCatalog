package filter

import (
	"errors"
)

var ErrInvalidOperator = errors.New("invalid operator")

const (
	DataTypeStr = "string"
	DataTypeInt = "int"

	OperatorEq               = "eq"
	OperatorNotEq            = "neq"
	OperatorLowerThan        = "lt"
	OperatorLowerThanEqual   = "lte"
	OperatorGreaterThan      = "gt"
	OperatorGreaterThanEqual = "gte"
	OperatorLike             = "like"
)

type Options struct {
	isToApply bool
	fields    []Field
}

func NewOptions() *Options {
	return &Options{}
}

type Field struct {
	Name     string
	Value    string
	Operator string
	Type     string
}

func (o *Options) IsToApply() bool {
	return o.isToApply
}

func (o *Options) AddField(name, operator, value, dType string) error {
	err := validateOperator(operator)
	if err != nil {
		return err
	}

	o.isToApply = true

	field := Field{
		Name:     name,
		Value:    value,
		Operator: operator,
		Type:     dType,
	}
	o.fields = append(o.fields, field)

	return nil
}

func (o *Options) Fields() []Field {
	return o.fields
}

func validateOperator(operator string) error {
	switch operator {
	case OperatorEq:
	case OperatorNotEq:
	case OperatorLowerThan:
	case OperatorLowerThanEqual:
	case OperatorGreaterThan:
	case OperatorGreaterThanEqual:
	case OperatorLike:
	default:
		return ErrInvalidOperator
	}
	return nil
}
