package migrations

import (
	h "spellslingerer.com/m/v2/migration_helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {

		err := h.LoadCsv("migration_files/sorin_signatures.csv", "cards", daos.New(db))
		if err != nil {
			return err
		}

		return h.LoadCsv("migration_files/sorin.csv", "spellslingers", daos.New(db))

	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
