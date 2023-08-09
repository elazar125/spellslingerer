package hooks

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func PreventSettingsChanges(app *pocketbase.PocketBase) {
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
}
