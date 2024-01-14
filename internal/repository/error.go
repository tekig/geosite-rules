package repository

import "errors"

var (
	ErrRuleNotFound      = errors.New("rule not found")
	ErrUnknownDomainType = errors.New("unknown domain type")
)
