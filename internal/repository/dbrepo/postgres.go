package dbrepo

import (
	"context"
	"time"

	"github.com/zepyrshut/breakfast-n-go/internal/models"
)

func (m *postgreDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *postgreDBRepo) InsertReservation(res models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `insert into reservations (first_name, second_name, last_name, 
		email, phone, arrival_date, departure_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.SecondName,
		res.LastName,
		res.Email,
		res.Phone,
		res.ArrivalDate,
		res.DepartureDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
