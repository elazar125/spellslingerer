package config

import (
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cast"
)

const COOKIE_NAME string = "pb_auth"

func LoadSettings(app *pocketbase.PocketBase) error {
	// load app settings from env variables
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		app.Settings().Smtp.Enabled = cast.ToBool(os.Getenv("SPELLSLINGERER_SMTP_ENABLED"))
		app.Settings().Smtp.Host = os.Getenv("SPELLSLINGERER_SMTP_HOST")
		app.Settings().Smtp.Port = cast.ToInt(os.Getenv("SPELLSLINGERER_SMTP_PORT"))
		app.Settings().Smtp.Username = os.Getenv("SPELLSLINGERER_SMTP_USERNAME")
		app.Settings().Smtp.Password = os.Getenv("SPELLSLINGERER_SMTP_PASSWORD")
		app.Settings().Meta.SenderName = os.Getenv("SPELLSLINGERER_SENDER_NAME")
		app.Settings().Meta.SenderAddress = os.Getenv("SPELLSLINGERER_SENDER_ADDRESS")

		app.Settings().Meta.AppName = "Spellslingerer"
		app.Settings().Meta.AppUrl = os.Getenv("SPELLSLINGERER_URL")
		app.Settings().Meta.HideControls = cast.ToBool(os.Getenv("SPELLSLINGERER_IS_PROD"))

		return nil
	})

	// prevent settings change
	app.OnSettingsBeforeUpdateRequest().Add(func(e *core.SettingsUpdateEvent) error {
		if e.OldSettings.Smtp.Enabled != e.NewSettings.Smtp.Enabled ||
			e.OldSettings.Smtp.Host != e.NewSettings.Smtp.Host ||
			e.OldSettings.Smtp.Port != e.NewSettings.Smtp.Port ||
			e.OldSettings.Smtp.Username != e.NewSettings.Smtp.Username ||
			e.OldSettings.Smtp.Password != e.NewSettings.Smtp.Password ||
			e.OldSettings.Meta.SenderName != e.NewSettings.Meta.SenderName ||
			e.OldSettings.Meta.SenderAddress != e.NewSettings.Meta.SenderAddress {
			return apis.NewForbiddenError("Cannot change the SMTP settings", nil)
		}

		if e.OldSettings.Meta.AppName != e.NewSettings.Meta.AppName ||
			e.OldSettings.Meta.AppUrl != e.NewSettings.Meta.AppUrl ||
			e.OldSettings.Meta.HideControls != e.NewSettings.Meta.HideControls {
			return apis.NewForbiddenError("Cannot change the app settings", nil)
		}

		return nil
	})

	return nil
}
