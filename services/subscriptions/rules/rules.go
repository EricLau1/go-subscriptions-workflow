package rules

import (
	"errors"
	"time"
)

const (
	DefaultPrice      = 50.0
	DefaultExpiration = time.Minute
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)