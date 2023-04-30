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
			"Special",
			"Deckbuilding",
			"Creature",
			"Spell",
			"Trap",
			"Artifact",
			"Land",
			"Skill",
		}

		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("types")
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
