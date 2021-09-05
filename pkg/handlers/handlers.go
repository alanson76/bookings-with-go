package handlers

import (
	"fmt"
	"net/http"

	"github.com/alanson76/playground/web/02_firstApp/pkg/config"
	"github.com/alanson76/playground/web/02_firstApp/pkg/models"
	"github.com/alanson76/playground/web/02_firstApp/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type or pattern
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

// Home is the handler for home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	fmt.Println(remoteIp)
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the handler for about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again!"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	fmt.Println("remote address got from session",remoteIp)
	stringMap["remote_ip"] = remoteIp

	// send the data to the template
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
