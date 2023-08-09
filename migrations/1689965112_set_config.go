package migrations

import (
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/spf13/cast"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		settings, _ := dao.FindSettings()
		settings.Smtp.Enabled = cast.ToBool(os.Getenv("SPELLSLINGERER_SMTP_ENABLED"))
		settings.Smtp.Host = os.Getenv("SPELLSLINGERER_SMTP_HOST")
		settings.Smtp.Port = cast.ToInt(os.Getenv("SPELLSLINGERER_SMTP_PORT"))
		settings.Smtp.Username = os.Getenv("SPELLSLINGERER_SMTP_USERNAME")
		settings.Smtp.Password = os.Getenv("SPELLSLINGERER_SMTP_PASSWORD")
		settings.Meta.SenderName = os.Getenv("SPELLSLINGERER_SENDER_NAME")
		settings.Meta.SenderAddress = os.Getenv("SPELLSLINGERER_SENDER_ADDRESS")

		settings.Meta.AppName = "Spellslingerer"
		settings.Meta.AppUrl = os.Getenv("SPELLSLINGERER_URL")
		settings.Meta.HideControls = cast.ToBool(os.Getenv("SPELLSLINGERER_IS_PROD"))

		dao.SaveSettings(settings, os.Getenv("SPELLSLINGERER_SETTINGS_ENCRYPTION_KEY"))

		return nil
	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
