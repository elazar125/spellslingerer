package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	h "spellslingerer.com/m/v2/migration_helpers"
)

func init() {
	m.Register(func(db dbx.Builder) error {

		dao := daos.New(db)

		err := h.LoadCsv("migration_files/cards_seismic_ids.csv", "cards", dao)
		if err != nil {
			return err
		}

		return h.LoadCsv("migration_files/spellslingers_seismic_ids.csv", "spellslingers", dao)

	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
