package shared

import (
	"errors"
	"time"
)

const (
	DefaultPrice      = 50.0
	DefaultExpiration = time.Second * 20
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)