package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zepyrshut/breakfast-n-go/internal/config"
	"github.com/zepyrshut/breakfast-n-go/internal/models"
	"github.com/zepyrshut/breakfast-n-go/internal/render"
)

// TemplateData holds data sent from handlers to template

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "hello again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Availability is the availability page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "availability.page.tmpl", &models.TemplateData{})
}

// Post Availability is the post availability page handler
func (m *Repository) DoPostAvailability(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("arrival")
	departure := r.Form.Get("departure")

	w.Write([]byte("arrival is" + arrival + " and departure is " + departure))
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
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Rooms is the rooms page handler
func (m *Repository) Rooms(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "rooms.page.tmpl", &models.TemplateData{})
}

// Terrace is the terrace page handler
func (m *Repository) Terrace(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "terrace.page.tmpl", &models.TemplateData{})
}

// Spa is the spa page handler
func (m *Repository) Spa(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "spa.page.tmpl", &models.TemplateData{})
}

// Reservation is the reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "reservation.page.tmpl", &models.TemplateData{})
}

// Contact is the contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
