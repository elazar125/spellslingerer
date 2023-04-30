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
			"Token",
			"Core",
			"Signature",
			"Common",
			"Rare",
			"Epic",
			"Mythic",
		}

		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("rarities")
		if err != nil {
			return err
		}

		for i, v := range values {
			model := models.NewRecord(collection)
			model.Set("name", v)
			model.Set("sort_order", i+1)

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
