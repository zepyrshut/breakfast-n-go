package models

import "time"

// User is the user model
type User struct {
	ID          int
	FirstName   string
	SecondName  string
	LastName    string
	Email       string
	Password    string
	AccessLevel string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservations model
type Reservation struct {
	ID            int
	FirstName     string
	SecondName    string
	LastName      string
	Email         string
	Phone         string
	ArrivalDate   time.Time
	DepartureDate time.Time
	RoomID        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Rooms         Room
}

// RoomRestriction is the room restrictions model
type RoomRestriction struct {
	ID            int
	ArrivalDate   time.Time
	DepartureDate time.Time
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
