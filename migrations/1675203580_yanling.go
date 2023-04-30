package migrations

import (
	h "spellslingerer.com/m/v2/migration_helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {

		err := h.LoadCsv("migration_files/yanling_signatures.csv", "cards", daos.New(db))
		if err != nil {
			return err
		}

		err = h.LoadCsv("migration_files/yanling.csv", "spellslingers", daos.New(db))
		if err != nil {
			return err
		}

		deck := h.Deck{
			Spellslinger: "Yanling",
			Land:         "Island",
			Cards: map[string]int{
				"Guard the Flock":   2,
				"Runeshell Crab":    2,
				"Filigree Fox":      2,
				"Welkin Tern":       2,
				"Unsummon":          2,
				"Stardust Moth":     2,
				"Yanling's Sparrow": 2,
				"Icy Manipulator":   2,
				"Divination":        1,
				"Call the Wind":     2,
				"Phantom Monster":   2,
				"Frost Lynx":        2,
				"Shell Game":        2,
				"River Dragon":      2,
				"Ejected!":          1,
				"Mahamoti Djinn":    2,
			},

			Code:  "DV1H4sIAAAAAAAAA23OIRKAMBBDUc9pkjaluxKPZOigev9boJgi4p77/zm285rj5gzwY1UxBMJRcuxwZDFsEY7qhshmKO7rtxoC6RiyTMOmMARXgvzxBa6LxjhXAQAA",
			Image: "/images/cards/tiles/Yanling's Sparrow.jpeg",
		}

		dao := daos.New(db)
		return h.LoadStarterDeck(dao, deck)

	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
