package main

import (
	"log"

	"spellslingerer.com/m/v2/authxtns"
	"spellslingerer.com/m/v2/config"
	"spellslingerer.com/m/v2/deck_hooks"
	"spellslingerer.com/m/v2/mailxtns"
	"spellslingerer.com/m/v2/router"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "spellslingerer.com/m/v2/migrations"
)

func main() {
	app := pocketbase.New()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := config.LoadSettings(app); err != nil {
		log.Fatal(err)
	}

	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{Automigrate: false})

	app.OnRecordAuthRequest().Add(authxtns.SetValidCookieHandler(app.Settings()))

	if err := mailxtns.BindMailEvents(app); err != nil {
		log.Fatal(err)
	}

	deck_hooks.BindDeckValidation(app)
	deck_hooks.BindDeckCodeHooks(app)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.HTTPErrorHandler = router.GetErrorHandler(app)

		e.Router.Renderer = &router.TemplateRegistry{
			Templates: router.LoadTemplates(),
		}

		e.Router.Pre(authxtns.LoadCookieContext(app))

		router.BindRoutes(e.Router, app.Dao(), app.Settings())

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
