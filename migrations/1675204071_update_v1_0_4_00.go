package migrations

import (
	"spellslingerer.com/m/v2/daoxtns"
	h "spellslingerer.com/m/v2/migration_helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {

		dao := daos.New(db)

		err := h.LoadCsv("migration_files/Update_v1.0.4.00_Cards.csv", "cards", dao)
		if err != nil {
			return err
		}

		jace, err := daoxtns.FindRecordByName(dao, "spellslingers", "Jace")
		if err != nil {
			return err
		}
		jace.Set("health", 27)
		err = dao.SaveRecord(jace)
		if err != nil {
			return err
		}

		serra, err := daoxtns.FindRecordByName(dao, "spellslingers", "Serra")
		if err != nil {
			return err
		}
		serra.Set("health", 20)
		return dao.SaveRecord(serra)

	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
