package daoxtns

import (
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func FindRecordByName(dao *daos.Dao, collection string, recordName string) (*models.Record, error) {
	return dao.FindFirstRecordByData(collection, "name", recordName)
}

func FindRecordBySeismicId(dao *daos.Dao, collection string, recordName string) (*models.Record, error) {
	return dao.FindFirstRecordByData(collection, "seismic_id", recordName)
}

func FindRecordById(
	dao *daos.Dao,
	collection string,
	id string,
	expandFields []string,
) (*models.Record, error) {
	record, err := dao.FindRecordById(collection, id)
	if err != nil {
		return nil, err
	}

	if expandFields != nil {
		expandFunc := func(c *models.Collection, ids []string) ([]*models.Record, error) {
			return dao.FindRecordsByIds(c.Name, ids)
		}

		failed := dao.ExpandRecord(record, expandFields, expandFunc)
		if len(failed) > 0 {
			return nil, errors.New(fmt.Sprintf("%v", failed))
		}
	}

	return record, nil
}

func FindRecordsByIds(
	dao *daos.Dao,
	collection string,
	ids []string,
	expandFields []string,
) ([]*models.Record, error) {
	records, err := dao.FindRecordsByIds(collection, ids)
	if err != nil {
		return nil, err
	}

	if expandFields != nil {
		expandFunc := func(c *models.Collection, ids []string) ([]*models.Record, error) {
			return dao.FindRecordsByIds(c.Name, ids)
		}

		failed := dao.ExpandRecords(records, expandFields, expandFunc)
		if len(failed) > 0 {
			return nil, errors.New(fmt.Sprintf("%v", failed))
		}
	}

	return records, nil
}

func FindAllRecords(
	dao *daos.Dao,
	collection string,
	expandFields []string,
) ([]*models.Record, error) {
	records, err := dao.FindRecordsByExpr(collection, dbx.NewExp("1 = 1"))
	if err != nil {
		return nil, err
	}

	if expandFields != nil {
		expandFunc := func(c *models.Collection, ids []string) ([]*models.Record, error) {
			return dao.FindRecordsByIds(c.Name, ids)
		}

		failed := dao.ExpandRecords(records, expandFields, expandFunc)
		if len(failed) > 0 {
			return nil, errors.New(fmt.Sprintf("%v", failed))
		}
	}

	return records, nil
}

func LoadDeckById(dao *daos.Dao, id string) (*models.Record, error) {
	expandFields := []string{
		"cards",
		"cards.set",
		"cards.colour",
		"cards.rarity",
		"cards.type",
		"land",
		"land.set",
		"spellslinger",
		"spellslinger.colour",
		"splash",
		"owner",
	}
	return FindRecordById(dao, "decks", id, expandFields)
}

func LoadCardById(dao *daos.Dao, id string) (*models.Record, error) {
	expandFields := []string{
		"set",
		"colour",
		"rarity",
		"type",
		"subtype",
		"reminders",
		"generates",
		"generates.set",
	}
	return FindRecordById(dao, "cards", id, expandFields)
}

func LoadMinimalCardById(dao *daos.Dao, id string) (*models.Record, error) {
	expandFields := []string{
		"set",
		"colour",
		"type",
		"subtype",
	}
	return FindRecordById(dao, "cards", id, expandFields)
}

func LoadMinimalCardsByIds(dao *daos.Dao, ids []string) ([]*models.Record, error) {
	expandFields := []string{
		"set",
		"colour",
		"type",
		"subtype",
	}
	return FindRecordsByIds(dao, "cards", ids, expandFields)
}

func LoadSpellslingers(dao *daos.Dao) ([]*models.Record, error) {
	expandFields := []string{
		"colour",
		"signatures",
		"signatures.set",
		"abilities",
	}
	return FindAllRecords(dao, "spellslingers", expandFields)
}

func LoadSpellslingerById(dao *daos.Dao, id string) (*models.Record, error) {
	expandFields := []string{
		"colour",
		"signatures",
		"signatures.set",
		"signatures.generates",
		"signatures.generates.set",
		"abilities",
		"abilities.generates",
		"abilities.generates.set",
		"starter_deck",
		"starter_deck.cards",
		"starter_deck.cards.colour",
		"starter_deck.cards.set",
		"starter_deck.land",
		"starter_deck.land.set",
		"starter_deck.spellslinger",
		"starter_deck.splash",
	}
	return FindRecordById(dao, "spellslingers", id, expandFields)
}

func LoadMinimalSpellslingerById(dao *daos.Dao, id string) (*models.Record, error) {
	return FindRecordById(dao, "spellslingers", id, nil)
}

func LoadUserById(dao *daos.Dao, id string) (*models.Record, error) {
	return dao.FindRecordById("users", id)
}
