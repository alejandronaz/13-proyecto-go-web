package services

import "errors"

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrProductExists           = errors.New("product already exists")
	ErrProductEmpty            = errors.New("product is empty")
	ErrInvalidExpirationFormat = errors.New("invalid expiration format")
)
