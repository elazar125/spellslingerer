package router

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
)

type TemplateRegistry struct {
	Templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return errors.New("Template not found -> " + name)
	}
	return tmpl.ExecuteTemplate(w, "page", data)
}

func LoadTemplates() map[string]*template.Template {
	funcs := template.FuncMap{
		"marshal": func(v any) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		"hasPrefix": func(str string, prefix string) bool {
			return strings.HasPrefix(str, prefix)
		},
		"getValue": func(record *models.Record, fieldName string) any {
			if val, ok := record.Expand()[fieldName]; ok {
				return val
			}
			return record.Get(fieldName)
		},
		"replace": func(str, from, to string) string {
			return strings.Replace(str, from, to, -1)
		},
		"split": func(str, at string) []string {
			return strings.Split(str, at)
		},
		"countCards": func(deck *models.Record) int {
			type cardDetail struct {
				Name             string `json:"name"`
				StandardQuantity int    `json:"standard_quantity"`
				FoilQuantity     int    `json:"foil_quantity"`
			}

			var cardDetails map[string]cardDetail
			cardDetailsJson := deck.GetString("card_details")
			err := json.Unmarshal([]byte(cardDetailsJson), &cardDetails)
			if err != nil {
				return 0
			}
			result := 0

			for _, detail := range cardDetails {
				result += detail.StandardQuantity
				result += detail.FoilQuantity
			}

			return result
		},
	}

	templates := make(map[string]*template.Template)

	templates["/"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/home.html", "templates/layout.html"))
	templates["/cards"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/cards.html", "templates/layout.html"))
	templates["/cards/:id"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/card.html", "templates/layout.html"))
	templates["/decks"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/decks.html", "templates/layout.html"))
	templates["/decks/:id"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/deck.html", "templates/layout.html"))
	templates["/decks/new"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/edit-deck.html", "templates/layout.html"))
	templates["/my-decks"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/decks.html", "templates/layout.html"))
	templates["/decks/:id/edit"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/edit-deck.html", "templates/layout.html"))
	templates["/users/:userid"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/decks.html", "templates/layout.html"))
	templates["/spellslingers"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/spellslingers.html", "templates/layout.html"))
	templates["/spellslingers/:id"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/spellslinger.html", "templates/layout.html"))
	templates["400"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/400.html", "templates/layout.html"))
	templates["401"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/401.html", "templates/layout.html"))
	templates["403"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/403.html", "templates/layout.html"))
	templates["404"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/404.html", "templates/layout.html"))
	templates["500"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/500.html", "templates/layout.html"))
	templates["/login"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/login.html", "templates/layout.html"))
	templates["/signup"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/signup.html", "templates/layout.html"))
	templates["/forgot-password"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/forgot-password.html", "templates/layout.html"))
	templates["/profile"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/profile.html", "templates/layout.html"))
	templates["/integrations"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/integrations.html", "templates/layout.html"))

	return templates
}
