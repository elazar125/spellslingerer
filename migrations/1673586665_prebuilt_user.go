package migrations

import (
	"os"

	h "spellslingerer.com/m/v2/migration_helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {

		return h.LoadUser(daos.New(db), "prebuiltdecks", "Prebuilt Decks", "spellslingererdeckbulider+prebuilt@gmail.com", os.Getenv("SPELLSLINGERER_PREBUILT_USER_PASSWORD"))

	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
