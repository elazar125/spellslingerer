package models

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"spellslingerer.com/m/v2/daoxtns"
)

type cardDetails struct {
	Name             string `json:"name"`
	StandardQuantity int    `json:"standard_quantity"`
	FoilQuantity     int    `json:"foil_quantity"`
}

const deckCodeIdentifier = "DV1"

func ImportDeckByCode(app core.App, code string, userId string) (string, error) {
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

	return loadDeck(app, cardQuantities, land, spellslinger, userId)
}

func decodeDeckCode(input string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf("%s.+$", deckCodeIdentifier))
	code := re.FindString(input)
	trimmed := strings.TrimPrefix(code, deckCodeIdentifier)
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

func loadDeck(app core.App, cards map[string]int, landId, spellslingerId, ownerId string) (string, error) {
	dao := app.Dao()

	spellslinger, err := daoxtns.FindRecordBySeismicId(dao, "spellslingers", spellslingerId)
	if err != nil {
		return "", err
	}

	land, err := daoxtns.FindRecordBySeismicId(dao, "cards", landId)
	if err != nil {
		return "", err
	}

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

	collection, err := dao.FindCollectionByNameOrId("decks")
	if err != nil {
		return "", err
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(app, record)
	form.LoadData(map[string]any{
		"image":        "/images/cards/tiles/default.jpeg",
		"name":         "Copied",
		"owner":        ownerId,
		"is_public":    false,
		"spellslinger": spellslinger.Id,
		"land":         land.Id,
		"cards":        cardIds,
		"card_details": string(cardDetailsJson),
	})
	err = form.Submit()
	if err != nil {
		return "", err
	}

	return record.Id, nil
}

func SetImportCode(dao *daos.Dao, deck *models.Record) error {
	fileContents, err := getFileContents(dao, deck)
	if err != nil {
		return err
	}

	code, err := encodeFileContents(fileContents)
	if err != nil {
		return err
	}

	deck.Set("code", fmt.Sprintf("%s%s", deckCodeIdentifier, code))

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

func ValidateDeck(dao *daos.Dao, deck *models.Record) validation.Errors {
	errs := make(validation.Errors, 0)
	cardIds := deck.GetStringSlice("cards")
	cardDetailsJson := deck.GetString("card_details")
	var cardDetails map[string]cardDetails
	err := json.Unmarshal([]byte(cardDetailsJson), &cardDetails)
	if err != nil {
		errs["load"] = validation.NewInternalError(err)
	}

	cards, err := daoxtns.LoadMinimalCardsByIds(dao, cardIds)
	if err != nil {
		errs["load"] = validation.NewInternalError(err)
	}

	land, err := daoxtns.LoadMinimalCardById(dao, deck.GetString("land"))
	if err != nil {
		errs["load"] = validation.NewInternalError(err)
	}

	spellslinger, err := daoxtns.LoadMinimalSpellslingerById(dao, deck.GetString("spellslinger"))
	if err != nil {
		errs["load"] = validation.NewInternalError(err)
	}

	var splash *models.Record
	if deck.GetString("splash") != "" {
		splash, err = dao.FindRecordById("colours", deck.GetString("splash"))
		if err != nil {
			errs["load"] = validation.NewInternalError(err)
		}
	}

	if !CanUseLand(spellslinger, land) {
		errs["land"] = validation.NewError("land", fmt.Sprintf("%s can not be used in this deck", land.GetString("name")))
	}

	totalCards := 0
	totalClasses := 0
	baseRuleExceptions := 0

	for _, card := range cards {
		details, exists := cardDetails[card.Id]
		cardName := card.GetString("name")
		if !exists {
			errs[cardName+"_quantity"] = validation.ErrMinGreaterEqualThanRequired.SetMessage(fmt.Sprintf("%s is missing a quantity", cardName))
		}

		full_quantity := details.StandardQuantity + details.FoilQuantity

		if details.StandardQuantity < 1 && details.FoilQuantity < 1 {
			errs[cardName+"_quantity"] = validation.ErrMinGreaterEqualThanRequired.SetMessage(fmt.Sprintf("Invalid quantity of %s", cardName))
		} else {
			totalCards += full_quantity
		}
		if card.GetBool("legendary") && details.StandardQuantity+details.FoilQuantity > 1 {
			errs[cardName+"_quantity"] = validation.ErrMaxLessEqualThanRequired.SetMessage(fmt.Sprintf("Too many copies of %s", cardName))
		}
		if !card.GetBool("legendary") && details.StandardQuantity+details.FoilQuantity > 2 {
			errs[cardName+"_quantity"] = validation.ErrMaxLessEqualThanRequired.SetMessage(fmt.Sprintf("Too many copies of %s", cardName))
		}
		if !IsValidSet(card, spellslinger) {
			errs[cardName+"_expansion"] = validation.NewError("validation_invalid_card", fmt.Sprintf("%s can not be used in this deck", cardName))
		} else if !SpellslingerCanUse(card, spellslinger, splash, full_quantity, &baseRuleExceptions) {
			errs[cardName+"_card"] = validation.NewError("validation_invalid_card", fmt.Sprintf("%s can not be used in this deck", cardName))
		}
		cardType := card.Expand()["type"].(*models.Record)
		subType := getSubTypes(card)
		if cardType.GetString("name") == "Skill" && len(subType) == 1 && subType[0].GetString("name") == "Class" {
			totalClasses++
		}
	}

	if totalClasses > 1 {
		errs["class"] = validation.ErrMaxLessEqualThanRequired.SetMessage(fmt.Sprintf("Too many classes in the deck, it has %d but maximum is 1", totalCards))
	}
	maxBaseRuleExceptions := maxBaseRuleExceptions(spellslinger)
	if baseRuleExceptions > maxBaseRuleExceptions {
		errs["cards"] = validation.ErrMaxLessEqualThanRequired.SetMessage(fmt.Sprintf("Too many special cards in deck, it has %d but maximum is %d", baseRuleExceptions, maxBaseRuleExceptions))
	}
	if totalCards > 30 {
		errs["total_count"] = validation.ErrMaxLessEqualThanRequired.SetMessage(fmt.Sprintf("Too many cards in the deck, it has %d but maximum is 30", totalCards))
	}

	return errs
}

func maxBaseRuleExceptions(spellslinger *models.Record) int {
	switch spellslinger.GetString("name") {
	case "Sorin":
		return 4
	case "Vivien":
		return 8
	case "Davriel":
		return 10
	default:
		return 6
	}
}

func SpellslingerCanUse(card, spellslinger, splash *models.Record, quantity int, baseRuleExceptions *int) bool {
	switch spellslinger.GetString("name") {
	case "Sorin":
		return SorinCanUse(card, spellslinger, splash, quantity, baseRuleExceptions)
	case "Davriel":
		return DavrielCanUse(card, spellslinger, quantity, baseRuleExceptions)
	case "Serra":
		return SerraCanUse(card)
	case "Vivien":
		return VivienCanUse(card, quantity, baseRuleExceptions)
	case "Drizzt":
		return DrizztCanUse(card, spellslinger, quantity, baseRuleExceptions)
	case "Nissa":
		return NissaCanUse(card, quantity, baseRuleExceptions)
	case "Chandra", "Jace", "Liliana", "Gideon", "Yanling":
		return MeetsColourOrSplash(card, spellslinger, splash, quantity, baseRuleExceptions)
	default:
		return MeetsColour(card, spellslinger)
	}
}

func SorinCanUse(card, spellslinger, splash *models.Record, quantity int, baseRuleExceptions *int) bool {
	if hasOnlyColour(card, "Black") || hasOnlyColour(card, "Colourless") {
		return true
	}

	if splash == nil {
		return false
	}

	cardColours := card.Expand()["colour"].([]*models.Record)
	ssColours := spellslinger.GetStringSlice("colour")

	for _, cardColour := range cardColours {
		ssHasColourOrSplash := false
		for _, ssColour := range append(ssColours, splash.Id) {
			if cardColour.Id == ssColour {
				ssHasColourOrSplash = true
			}
		}
		if !ssHasColourOrSplash {
			return false
		}
	}

	cardType := card.Expand()["type"].(*models.Record)
	subTypes := getSubTypes(card)
	isVampire := false
	if cardType.GetString("name") == "Creature" {
		for _, subType := range subTypes {
			if subType.GetString("name") == "Vampire" {
				isVampire = true
			}
		}
	}
	if !isVampire {
		*baseRuleExceptions += quantity
	}

	return true
}

func DavrielCanUse(card, spellslinger *models.Record, quantity int, baseRuleExceptions *int) bool {
	if hasOnlyColour(card, "Black") || hasOnlyColour(card, "Colourless") {
		return true
	}

	if quantity != 1 {
		return false
	}

	*baseRuleExceptions += quantity
	return true
}

func SerraCanUse(card *models.Record) bool {
	if hasOnlyColour(card, "Colourless") {
		return true
	}

	cardColours := card.Expand()["colour"].([]*models.Record)
	for _, colour := range cardColours {
		if colour.GetString("name") == "White" {
			return true
		}
	}

	return false
}

func VivienCanUse(card *models.Record, quantity int, baseRuleExceptions *int) bool {
	if hasOnlyColour(card, "Green") || hasOnlyColour(card, "Colourless") {
		return true
	}

	cardType := card.Expand()["type"].(*models.Record)
	if cardType.GetString("name") == "Creature" {
		*baseRuleExceptions += quantity
		return true
	}

	return false
}

func DrizztCanUse(card, spellslinger *models.Record, quantity int, baseRuleExceptions *int) bool {
	if card.GetBool("legendary") {
		*baseRuleExceptions += quantity
		return true
	}

	ssColours := spellslinger.GetStringSlice("colour")
	cardColours := card.GetStringSlice("colour")

	for _, cardColour := range cardColours {
		ssHasColour := false
		for _, ssColour := range ssColours {
			if cardColour == ssColour {
				ssHasColour = true
			}
		}
		if !ssHasColour {
			return false
		}
	}
	return true
}

func NissaCanUse(card *models.Record, quantity int, baseRuleExceptions *int) bool {
	if hasOnlyColour(card, "Green") || hasOnlyColour(card, "Colourless") {
		return true
	}

	*baseRuleExceptions += quantity
	return true
}

func MeetsColourOrSplash(card, spellslinger, splash *models.Record, quantity int, baseRuleExceptions *int) bool {
	if hasOnlyColour(card, "Colourless") {
		return true
	}

	ssColours := spellslinger.GetStringSlice("colour")
	cardColours := card.GetStringSlice("colour")

	for _, cardColour := range cardColours {
		ssHasColour := false
		for _, ssColour := range ssColours {
			if cardColour == ssColour {
				ssHasColour = true
			}
		}
		if !ssHasColour {
			if splash == nil {
				return false
			}
			ssHasColourOrSplash := false
			for _, ssColour := range append(ssColours, splash.Id) {
				if cardColour == ssColour {
					ssHasColourOrSplash = true
				}
			}
			if !ssHasColourOrSplash {
				return false
			} else {
				*baseRuleExceptions += quantity
			}
		}
	}
	return true
}

func MeetsColour(card, spellslinger *models.Record) bool {
	if hasOnlyColour(card, "Colourless") {
		return true
	}

	ssColours := spellslinger.GetStringSlice("colour")
	cardColours := card.GetStringSlice("colour")

	for _, cardColour := range cardColours {
		ssHasColour := false
		for _, ssColour := range ssColours {
			if cardColour == ssColour {
				ssHasColour = true
			}
		}
		if !ssHasColour {
			return false
		}
	}
	return true
}

func hasOnlyColour(card *models.Record, colour string) bool {
	colours := card.Expand()["colour"].([]*models.Record)
	if len(colours) == 1 && colours[0].GetString("name") == colour {
		return true
	}
	return false
}

func IsValidSet(card, spellslinger *models.Record) bool {
	set := card.Expand()["set"].(*models.Record)
	if set.GetString("name") != "Signatures" {
		return true
	}
	for _, id := range spellslinger.GetStringSlice("signatures") {
		if card.Id == id {
			return true
		}
	}
	return false
}

func CanUseLand(spellslinger, land *models.Record) bool {
	if spellslinger.GetString("name") == "Nissa" {
		return true
	}

	if hasOnlyColour(land, "Colourless") {
		return true
	}

	ssColours := spellslinger.GetStringSlice("colour")
	landColours := land.GetStringSlice("colour")

	for _, landColour := range landColours {
		ssHasColour := false
		for _, ssColour := range ssColours {
			if landColour == ssColour {
				ssHasColour = true
			}
		}
		if !ssHasColour {
			return false
		}
	}
	return true
}

func getSubTypes(card *models.Record) []*models.Record {
	if card.Expand()["subtype"] != nil {
		return card.Expand()["subtype"].([]*models.Record)
	}
	return []*models.Record{}
}
