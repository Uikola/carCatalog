package errorz

import "errors"

var ErrNoCars = errors.New("no cars")
var ErrCarNotFound = errors.New("car not found")
var ErrInvalidRegNum = errors.New("invalid registration number")
var ErrCarAlreadyExists = errors.New("car already exists")
