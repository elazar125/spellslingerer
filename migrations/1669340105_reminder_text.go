package migrations

import (
	h "spellslingerer.com/m/v2/migration_helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {

		return h.LoadCsv("migration_files/reminder_text.csv", "reminders", daos.New(db))

	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
