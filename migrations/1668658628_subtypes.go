package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		values := []string{
			"???",
			"Angel",
			"Artificer",
			"Azra",
			"Barbarian",
			"Bard",
			"Basic",
			"Beast",
			"Bird",
			"Blessing",
			"Cat",
			"Civilian",
			"Class",
			"Cleric",
			"Construct",
			"Deal",
			"Demon",
			"Devil",
			"Dinosaur",
			"Dog",
			"Dragon",
			"Drow",
			"Dwarf",
			"Elemental",
			"Elf",
			"Fighter",
			"Giant",
			"Goblin",
			"God",
			"Homunculus",
			"Illusion",
			"Leviathan",
			"Merfolk",
			"Monk",
			"Monster",
			"Mystic",
			"Ooze",
			"Passive",
			"Pest",
			"Pirate",
			"Plant",
			"Ranger",
			"Rogue",
			"Snake",
			"Soldier",
			"Sorcerer",
			"Spirit",
			"Vampire",
			"Warrior",
			"Werewolf",
			"Wizard",
			"Zombie",
		}

		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("subtypes")
		if err != nil {
			return err
		}

		for _, v := range values {
			model := models.NewRecord(collection)
			model.Set("name", v)

			err = dao.SaveRecord(model)
			if err != nil {
				return err
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
