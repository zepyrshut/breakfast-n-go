package dbrepo

import (
	"context"
	"time"

	"github.com/zepyrshut/breakfast-n-go/internal/models"
)

func (m *postgreDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgreDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into reservations (first_name, second_name, last_name,
		email, phone, arrival_date, departure_date, room_id, created_at,
		updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
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
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction insert a restriction into the database
func (m *postgreDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (arrival_date, departure_date,
		room_id, reservation_id, restriction_id, created_at, updated_at) values
		($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.ArrivalDate,
		r.DepartureDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesAndRoomID return true if availability exist, and false if no
func (m *postgreDBRepo) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `select count(id) from room_restrictions rr where room_id = $1 and $2 <
	departure_date and $3 > arrival_date `

	var numRows int

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns all rooms are availables for given date range
func (m *postgreDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select r.id, r.room_name from rooms r where r.id not in (select
		rr.room_id from room_restrictions rr where $1 < departure_date
		and $2 > arrival_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}
