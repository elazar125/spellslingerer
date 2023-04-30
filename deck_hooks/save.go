package deck_hooks

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"spellslingerer.com/m/v2/daoxtns"
	"strings"
)

func BindDeckCodeHooks(app *pocketbase.PocketBase) {
	app.OnRecordBeforeCreateRequest("decks").Add(func(e *core.RecordCreateEvent) error {
		return setImportCode(app.Dao(), e.Record)
	})
	app.OnRecordBeforeUpdateRequest("decks").Add(func(e *core.RecordUpdateEvent) error {
		return setImportCode(app.Dao(), e.Record)
	})
}

func setImportCode(dao *daos.Dao, deck *models.Record) error {
	fileContents, err := getFileContents(dao, deck)
	if err != nil {
		return err
	}

	code, err := encodeFileContents(fileContents)
	if err != nil {
		return err
	}

	deck.Set("code", fmt.Sprintf("DV1%s", code))

	return nil
}

func encodeFileContents(contents []byte) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(contents); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func getFileContents(dao *daos.Dao, deck *models.Record) ([]byte, error) {
	ids := []string{}

	cardIds := deck.GetStringSlice("cards")
	cardDetailsJson := deck.GetString("card_details")
	var cardDetails map[string]cardDetails
	err := json.Unmarshal([]byte(cardDetailsJson), &cardDetails)
	if err != nil {
		return []byte{}, err
	}

	cards, err := daoxtns.LoadMinimalCardsByIds(dao, cardIds)
	if err != nil {
		return []byte{}, err
	}

	land, err := daoxtns.LoadMinimalCardById(dao, deck.GetString("land"))
	if err != nil {
		return []byte{}, err
	}

	spellslinger, err := daoxtns.LoadMinimalSpellslingerById(dao, deck.GetString("spellslinger"))
	if err != nil {
		return []byte{}, err
	}

	ids = append(ids, spellslinger.GetString("seismic_id"))
	ids = append(ids, land.GetString("seismic_id"))

	for _, card := range cards {
		details, exists := cardDetails[card.Id]
		if !exists {
			return []byte{}, errors.New("A card was missing from the Card Details JSON")
		}
		for i := 0; i < details.FoilQuantity+details.StandardQuantity; i++ {
			ids = append(ids, card.GetString("seismic_id"))
		}
	}

	return []byte(strings.Join(ids, "\n")), nil
}
