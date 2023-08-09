package hooks

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"spellslingerer.com/m/v2/models"
)

func BindDeckHooks(app *pocketbase.PocketBase) {
	app.OnRecordBeforeCreateRequest("decks").Add(func(e *core.RecordCreateEvent) error {
		err := models.SetImportCode(app.Dao(), e.Record)
		if err != nil {
			return err
		}
		validationErrors := models.ValidateDeck(app.Dao(), e.Record)
		if len(validationErrors) > 0 {
			return validationErrors
		}
		return nil
	})
	app.OnRecordBeforeUpdateRequest("decks").Add(func(e *core.RecordUpdateEvent) error {
		err := models.SetImportCode(app.Dao(), e.Record)
		if err != nil {
			return err
		}
		validationErrors := models.ValidateDeck(app.Dao(), e.Record)
		if len(validationErrors) > 0 {
			return validationErrors
		}
		return nil
	})
}
