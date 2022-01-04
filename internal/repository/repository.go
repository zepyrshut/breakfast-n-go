package repository

import "github.com/zepyrshut/breakfast-n-go/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
