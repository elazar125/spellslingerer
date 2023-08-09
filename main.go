package main

import (
	"log"

	"spellslingerer.com/m/v2/authxtns"
	"spellslingerer.com/m/v2/hooks"
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

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{Automigrate: false})

	hooks.PreventSettingsChanges(app)

	app.OnRecordAuthRequest().Add(authxtns.SetValidCookieHandler(app.Settings()))

	mailxtns.BindMailEvents(app)
	hooks.BindDeckHooks(app)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.HTTPErrorHandler = router.GetErrorHandler(e.App)

		e.Router.Renderer = &router.TemplateRegistry{
			Templates: router.LoadTemplates(),
		}

		e.Router.Pre(authxtns.LoadCookieContext(e.App))

		router.BindRoutes(e.App, e.Router)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
