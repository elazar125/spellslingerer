package models

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"io"
	"spellslingerer.com/m/v2/daoxtns"
	"strings"
)

type cardDetails struct {
	Name             string `json:"name"`
	StandardQuantity int    `json:"standard_quantity"`
	FoilQuantity     int    `json:"foil_quantity"`
}

func ImportDeckByCode(dao *daos.Dao, code string) (string, error) {
	deckIds, err := decodeDeckCode(code)
	if err != nil {
		return "", err
	}
	lines := strings.Split(deckIds, "\n")
	spellslinger := lines[0]
	land := lines[1]
	cards := lines[2:]
	cardQuantities := make(map[string]int)
	for _, card := range cards {
		cardQuantities[card]++
	}

	return loadDeck(dao, cardQuantities, land, spellslinger, code, "")
}

func decodeDeckCode(code string) (string, error) {
	trimmed := strings.TrimPrefix(code, "DV1")
	gzipped, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return "", err
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(gzipped))
	if err != nil {
		return "", err
	}
	defer gz.Close()

	data, err := io.ReadAll(gz)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func loadDeck(dao *daos.Dao, cards map[string]int, landId, spellslingerId, code, ownerId string) (string, error) {
	collection, err := dao.FindCollectionByNameOrId("decks")
	if err != nil {
		return "", err
	}

	record := models.NewRecord(collection)
	record.Set("image", "/images/cards/tiles/default.jpeg")
	record.Set("code", code)
	record.Set("name", "Copied")
	record.Set("owner", ownerId)
	record.Set("is_public", false)

	spellslinger, err := daoxtns.FindRecordBySeismicId(dao, "spellslingers", spellslingerId)
	if err != nil {
		return "", err
	}
	record.Set("spellslinger", spellslinger.Id)

	land, err := daoxtns.FindRecordBySeismicId(dao, "cards", landId)
	if err != nil {
		return "", err
	}
	record.Set("land", land.Id)

	cardIds := make([]string, 0)
	cardDetails := make(map[string]map[string]int, 0)

	for cardId, count := range cards {
		card, err := daoxtns.FindRecordBySeismicId(dao, "cards", cardId)
		if err != nil {
			return "", err
		}
		cardIds = append(cardIds, card.Id)
		cardDetails[card.Id] = map[string]int{
			"standard_quantity": count,
			"foil_quantity":     0,
		}
	}

	cardDetailsJson, err := json.Marshal(cardDetails)
	if err != nil {
		return "", err
	}

	record.Set("cards", cardIds)
	record.Set("card_details", cardDetailsJson)

	// TODO: Use form for validation
	err = dao.SaveRecord(record)
	if err != nil {
		return "", err
	}

	return record.Id, nil
}
