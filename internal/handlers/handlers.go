package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/zepyrshut/breakfast-n-go/internal/config"
	"github.com/zepyrshut/breakfast-n-go/internal/driver"
	"github.com/zepyrshut/breakfast-n-go/internal/forms"
	"github.com/zepyrshut/breakfast-n-go/internal/helpers"
	"github.com/zepyrshut/breakfast-n-go/internal/models"
	"github.com/zepyrshut/breakfast-n-go/internal/render"
	"github.com/zepyrshut/breakfast-n-go/internal/repository"
	"github.com/zepyrshut/breakfast-n-go/internal/repository/dbrepo"
)

// TemplateData holds data sent from handlers to template

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Availability is the availability page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "availability.page.tmpl", &models.TemplateData{})
}

// Post Availability is the post availability page handler
func (m *Repository) DoPostAvailability(w http.ResponseWriter, r *http.Request) {

	ad := r.Form.Get("arrival_date")
	dd := r.Form.Get("departure_date")

	layout := "02/01/2006"
	arrivalDate, err := time.Parse(layout, ad)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	departureDate, err := time.Parse(layout, dd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(arrivalDate, departureDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No hay habitaciones disponibles")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		ArrivalDate:   arrivalDate,
		DepartureDate: departureDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	w.Write([]byte("arrival is " + ad + " and departure is " + dd))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK:      true,
		Message: "available!",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Rooms is the rooms page handler
func (m *Repository) Rooms(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "rooms.page.tmpl", &models.TemplateData{})
}

// Terrace is the terrace page handler
func (m *Repository) Terrace(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "terrace.page.tmpl", &models.TemplateData{})
}

// Spa is the spa page handler
func (m *Repository) Spa(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "spa.page.tmpl", &models.TemplateData{})
}

// Reservation is the reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	ad := r.Form.Get("arrival_date")
	dd := r.Form.Get("departure_date")

	// 2022-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "02/01/2006"
	// TODO: beware for date format in form
	arrivalDate, err := time.Parse(layout, ad)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	departureDate, err := time.Parse(layout, dd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName:     r.Form.Get("first_name"),
		SecondName:    r.Form.Get("second_name"),
		LastName:      r.Form.Get("last_name"),
		Email:         r.Form.Get("email"),
		Phone:         r.Form.Get("phone"),
		ArrivalDate:   arrivalDate,
		DepartureDate: departureDate,
		RoomID:        roomID,
	}

	form := forms.New(r.PostForm)

	// Form validation checkers
	//form.Required("first_name", "second_name", "email")
	//form.Required("first_name")
	//form.MinLength("first_name", 3, r)
	//form.Has("phone", r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		ArrivalDate:   arrivalDate,
		DepartureDate: departureDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Reservation summary page handler
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("can't get error from session")
		m.App.Session.Put(r.Context(), "error", "No se pudo obtener la reserva")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Contact is the contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
