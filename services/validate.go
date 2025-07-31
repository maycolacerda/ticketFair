package services

import (
	"github.com/maycolacerda/ticketfair/models"
)

type Validator interface {
	Validate(models.User) []string
}
