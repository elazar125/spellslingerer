package migrations

import (
	h "spellslingerer.com/m/v2/migration_helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {

		dao := daos.New(db)

		err := h.LoadCsv("migration_files/Update_v1.0.5.00_Signatures.csv", "cards", dao)
		if err != nil {
			return err
		}

		err = h.LoadCsv("migration_files/Update_v1.0.5.00_Spellslingers.csv", "spellslingers", dao)
		if err != nil {
			return err
		}

		garruk := h.Deck{
			Spellslinger: "Garruk",
			Land:         "Bayou",
			Cards: map[string]int{
				"Gladecover Scout":     2,
				"Kalonian Tusker":      2,
				"Elvish Infuser":       2,
				"Grudge Match":         2,
				"Cursebearer":          2,
				"Rotting Baboon":       2,
				"Courtly Killer":       2,
				"Spiked Baloth":        2,
				"Wildspeaker's Fury":   2,
				"Flagrant Foul":        2,
				"Peafoul":              2,
				"Ornery Leotau":        2,
				"Briarhorn":            2,
				"Veil-Cursed Predator": 2,
				"Colossal Dreadmaw":    2,
			},
			Code:  "DV1H4sIAAAAAAAAA23PsQ7CMBAD0J1/ieRzksYZmViYKkTHDm35/08AFCGq021vsM6+wlc1VCU09lSakIStp53HblNj1YbL/bEuT1tl/BGtRxQicoqI7Hi7zh/S8Rsws4DIjCjPcSw7jukMSPSA9q84s8hxVBTH8bxOfANFBZeeeQEAAA==",
			Image: "/images/cards/tiles/Spiked Baloth.jpeg",
		}

		davriel := h.Deck{
			Spellslinger: "Davriel",
			Land:         "Swamp",
			Cards: map[string]int{
				"Decaying Ghoul":    2,
				"Shock":             1,
				"Filigree Fox":      2,
				"Mend":              1,
				"Unsummon":          1,
				"Arcane Flight":     1,
				"Drain Blood":       2,
				"Daggerclaw Imp":    2,
				"Famished Fiend":    2,
				"Night's Whisper":   2,
				"Forest Patrol":     1,
				"Paincast Demon":    2,
				"Courtly Killer":    2,
				"Davriel's Scourge": 2,
				"Trumpet Blast":     1,
				"Thwack!":           1,
				"Shell Game":        1,
				"Deathspeaker Naga": 2,
				"Bloodlord":         2,
			},
			Code:  "DV1H4sIAAAAAAAAA23PIQ7AMAwDQL7XxGnSpnDS4OC0wf7/F0NtCsxOBrZ8ncf9jO/FCNFJb5VQ1BaN0QpWilVm8EktSliFEWB0BKFay5RRcm2joxNqrEPF82b0ZGXsO383Cqn0VwEAAA==",
			Image: "/images/cards/tiles/Night's Whisper.jpeg",
		}

		err = h.LoadStarterDeck(dao, garruk)
		if err != nil {
			return err
		}
		return h.LoadStarterDeck(dao, davriel)

	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
