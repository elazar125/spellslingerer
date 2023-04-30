package migration_helpers

import (
	"fmt"

	"spellslingerer.com/m/v2/daoxtns"

	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Deck struct {
	Spellslinger string
	Land         string
	Cards        map[string]int
	Code         string
	Image        string
}

func LoadCsv(filePath string, collectionName string, dao *daos.Dao) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	collection, err := dao.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return err
	}

	r := csv.NewReader(file)

	header, err := r.Read()
	if err != nil {
		return err
	}

	importFuncMap, err := buildImportFuncMap(header)
	if err != nil {
		return err
	}

	lines, err := r.ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		record, err := getMatchingRecord(dao, collection, header, line)
		if err != nil {
			return err
		}

		for i, v := range line {
			err = importFuncMap[i](dao, v, record)
			if err != nil {
				return err
			}
		}

		err = dao.SaveRecord(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadUser(dao *daos.Dao, username, displayName, email, password string) error {
	collection, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	record.SetUsername(username)
	record.Set("display_name", displayName)
	record.SetEmail(email)
	record.SetEmailVisibility(false)
	record.SetPassword(password)
	record.SetVerified(true)

	return dao.SaveRecord(record)
}

func LoadStarterDeck(dao *daos.Dao, deck Deck) error {
	prebuiltUser, err := dao.FindAuthRecordByUsername("users", "prebuiltdecks")
	if err != nil {
		return err
	}

	description := fmt.Sprintf("The starting deck made available when %s is unlocked", deck.Spellslinger)
	deckName := fmt.Sprintf("Starter: %s", deck.Spellslinger)

	deckId, err := LoadDeck(dao, deck.Cards, deck.Land, deck.Spellslinger, deck.Image, deck.Code, description, deckName, prebuiltUser.Id)
	if err != nil {
		return err
	}

	spellslinger, err := daoxtns.FindRecordByName(dao, "spellslingers", deck.Spellslinger)
	spellslinger.Set("starter_deck", deckId)
	return dao.SaveRecord(spellslinger)
}

func LoadDeck(dao *daos.Dao, cards map[string]int, landName, spellslingerName, tileImage, code, description, name, ownerId string) (string, error) {
	collection, err := dao.FindCollectionByNameOrId("decks")
	if err != nil {
		return "", err
	}

	record := models.NewRecord(collection)
	record.Set("image", tileImage)
	record.Set("code", code)
	record.Set("description", description)
	record.Set("name", name)
	record.Set("owner", ownerId)
	record.Set("is_public", true)

	spellslinger, err := daoxtns.FindRecordByName(dao, "spellslingers", spellslingerName)
	if err != nil {
		return "", err
	}
	record.Set("spellslinger", spellslinger.Id)

	land, err := daoxtns.FindRecordByName(dao, "cards", landName)
	if err != nil {
		return "", err
	}
	record.Set("land", land.Id)

	cardIds := make([]string, 0)
	cardDetails := make(map[string]map[string]int, 0)

	for cardName, count := range cards {
		card, err := daoxtns.FindRecordByName(dao, "cards", cardName)
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

	err = dao.SaveRecord(record)
	if err != nil {
		return "", err
	}

	return record.Id, nil
}

func buildImportFuncMap(header []string) ([]func(*daos.Dao, string, *models.Record) error, error) {
	importFuncMap := make([]func(*daos.Dao, string, *models.Record) error, len(header))

	for i, v := range header {
		importFunc, err := mapHeaderToImportFunc(v)
		if err != nil {
			return nil, err
		}
		importFuncMap[i] = importFunc
	}
	return importFuncMap, nil
}

func mapHeaderToImportFunc(header string) (func(*daos.Dao, string, *models.Record) error, error) {
	switch header {
	case "Rarity":
		return setRarityOnRecord, nil
	case "Color":
		return setColorOnRecord, nil
	case "Cost":
		return setCostOnRecord, nil
	case "Title":
		return setTitleOnRecord, nil
	case "Power":
		return setPowerOnRecord, nil
	case "Health":
		return setHealthOnRecord, nil
	case "Type":
		return setTypeOnRecord, nil
	case "SubType":
		return setSubTypeOnRecord, nil
	case "Text":
		return setCardTextOnRecord, nil
	case "Exp":
		return setExpOnRecord, nil
	case "Legendary":
		return setLegendaryOnRecord, nil
	case "Land Average Percentage to Hit":
		return setChanceOnRecord, nil
	case "Charges":
		return setChargesOnRecord, nil
	case "Artist":
		return setArtistOnRecord, nil
	case "Generates":
		return setGeneratesOnRecord, nil
	case "Reminders":
		return setRemindersOnRecord, nil
	case "Name":
		return setNameOnRecord, nil
	case "Reminder Text":
		return setReminderTextOnRecord, nil
	case "Ability":
		return setSpellslingerAbilityOnRecord, nil
	case "Signatures":
		return setSpellslingerSignaturesOnRecord, nil
	case "Seismic ID":
		return setSeismicIDOnRecord, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unexpected column found in migration CSV: %s", header))
	}
}

func setRarityOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("rarity", value, record) {
		return nil
	}
	rarityMap := map[string]string{
		"T": "Token",
		"B": "Core",
		"S": "Signature",
		"C": "Common",
		"R": "Rare",
		"E": "Epic",
		"M": "Mythic",
	}
	rarity, err := daoxtns.FindRecordByName(dao, "rarities", rarityMap[value])
	if err != nil {
		return err
	}
	record.Set("rarity", rarity.Id)
	return nil
}

func setColorOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("colour", value, record) {
		return nil
	}
	colourMap := map[string]string{
		"W": "White",
		"U": "Blue",
		"B": "Black",
		"R": "Red",
		"G": "Green",
		"C": "Colourless",
	}
	colourIds := make([]string, 0)
	for _, c := range []rune(value) {
		colour, err := daoxtns.FindRecordByName(dao, "colours", colourMap[string(c)])
		if err != nil {
			return err
		}
		colourIds = append(colourIds, colour.Id)
	}
	record.Set("colour", colourIds)
	return nil
}

func setCostOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("cost", value, record) {
		return nil
	}
	cost, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return err
	}
	record.Set("cost", cost)
	return nil
}

func setTitleOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("name", value, record) {
		return nil
	}
	record.Set("name", value)
	return nil
}

func setPowerOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("power", value, record) {
		return nil
	}
	power, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return err
	}
	record.Set("power", power)
	return nil
}

func setHealthOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("health", value, record) {
		return nil
	}
	health, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return err
	}
	record.Set("health", health)
	return nil
}

func setTypeOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("type", value, record) {
		return nil
	}
	cardType, err := daoxtns.FindRecordByName(dao, "types", value)
	if err != nil {
		return err
	}
	record.Set("type", cardType.Id)
	return nil
}

func setSubTypeOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("subtype", value, record) {
		return nil
	}
	subtypes := strings.Fields(value)
	stIds := make([]string, 0)
	for _, v := range subtypes {
		st, err := daoxtns.FindRecordByName(dao, "subtypes", v)
		if err != nil {
			return err
		}
		stIds = append(stIds, st.Id)
	}
	record.Set("subtype", stIds)
	return nil
}

func setCardTextOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("ability", value, record) {
		return nil
	}
	record.Set("ability", value)
	return nil
}

func setExpOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("set", value, record) {
		return nil
	}
	var setName string
	if len(value) == 3 {
		setName = "Signatures"
	} else {
		setName = value
	}
	set, err := daoxtns.FindRecordByName(dao, "sets", setName)
	if err != nil {
		return err
	}
	record.Set("set", set.Id)
	return nil
}

func setLegendaryOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("legendary", value, record) {
		return nil
	}
	record.Set("legendary", value == "Y")
	return nil
}

func setChanceOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("chance", value, record) {
		return nil
	}
	chance, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return err
	}
	record.Set("chance", chance)
	return nil
}

func setChargesOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("charges", value, record) {
		return nil
	}
	charges, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return err
	}
	record.Set("charges", charges)
	return nil
}

func setArtistOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("artist", value, record) {
		return nil
	}
	record.Set("artist", value)
	return nil
}

func setGeneratesOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("generates", value, record) {
		return nil
	}
	cards := strings.Split(value, ";")
	cIds := make([]string, 0)
	for _, v := range cards {
		c, err := daoxtns.FindRecordByName(dao, "cards", v)
		if err != nil {
			return err
		}
		cIds = append(cIds, c.Id)
	}
	record.Set("generates", cIds)
	return nil
}

func setRemindersOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("reminders", value, record) {
		return nil
	}
	reminders := strings.Split(value, ";")
	rIds := make([]string, 0)
	for _, v := range reminders {
		r, err := daoxtns.FindRecordByName(dao, "reminders", v)
		if err != nil {
			return err
		}
		rIds = append(rIds, r.Id)
	}
	record.Set("reminders", rIds)
	return nil
}

func setNameOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("name", value, record) {
		return nil
	}
	record.Set("name", value)
	return nil
}

func setReminderTextOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("text", value, record) {
		return nil
	}
	record.Set("text", value)
	return nil
}

func setSpellslingerAbilityOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("abilities", value, record) {
		return nil
	}
	cards := strings.Split(value, ";")
	cIds := make([]string, 0)
	for _, v := range cards {
		c, err := daoxtns.FindRecordByName(dao, "cards", v)
		if err != nil {
			return err
		}
		cIds = append(cIds, c.Id)
	}
	record.Set("abilities", cIds)
	return nil
}

func setSpellslingerSignaturesOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("signatures", value, record) {
		return nil
	}
	cards := strings.Split(value, ";")
	stIds := make([]string, 0)
	for _, v := range cards {
		st, err := daoxtns.FindRecordByName(dao, "cards", v)
		if err != nil {
			return err
		}
		stIds = append(stIds, st.Id)
	}
	record.Set("signatures", stIds)
	return nil
}

func setSeismicIDOnRecord(dao *daos.Dao, value string, record *models.Record) error {
	if setDefaultValue("seismic_id", value, record) {
		return nil
	}
	record.Set("seismic_id", value)
	return nil
}

func setDefaultValue(fieldName string, value string, record *models.Record) bool {
	if value != "" {
		return false
	}
	if record.Get(fieldName) == nil {
		record.Set(fieldName, "")
	}
	return true
}

func getMatchingRecord(dao *daos.Dao, collection *models.Collection, header, line []string) (*models.Record, error) {
	findIndex := func(toFind string, strs []string) (int, error) {
		for i, str := range strs {
			if str == toFind {
				return i, nil
			}
		}
		return -1, errors.New("Identifying column not found")
	}

	var columnName string

	switch collection.Name {
	case "reminders", "spellslingers":
		columnName = "Name"
		break
	case "cards":
		columnName = "Title"
		break
	}

	index, err := findIndex(columnName, header)
	if err != nil {
		return nil, err
	}
	record, err := daoxtns.FindRecordByName(dao, collection.Name, line[index])
	if err != nil {
		return models.NewRecord(collection), nil
	}
	return record, nil
}
