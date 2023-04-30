package router

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"spellslingerer.com/m/v2/daoxtns"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/settings"
)

type ViewData struct {
	Title       string
	Url         string
	Description string
	Image       string
	ImageAlt    string
	Domain      string
	Content     any
}

type DeckListPageData struct {
	Title  string
	Filter template.JS
}

type EditDeckPageData struct {
	Deck          *models.Record
	Spellslingers []*models.Record
	Colours       []*models.Record
}

type DeckPageData struct {
	Deck        *models.Record
	CurrentUser *models.Record
}

func DefaultViewData(c echo.Context, settings *settings.Settings) ViewData {
	return ViewData{
		Title:       "Spellslingerer",
		Description: "A card database and deckbuilder companion to MtG Spellslingers",
		Url:         c.Path(),
		Image:       "/images/Spellslingerer.webp",
		ImageAlt:    "Spellslingerer.com logo",
		Domain:      settings.Meta.AppUrl,
		Content:     nil,
	}
}

func LoadDeck(c echo.Context, dao *daos.Dao, settings *settings.Settings) (ViewData, error) {
	userIntfc := c.Get(apis.ContextAuthRecordKey)
	var user *models.Record
	if userIntfc != nil {
		user = userIntfc.(*models.Record)
	} else {
		users, err := dao.FindCollectionByNameOrId("users")
		if err != nil {
			return ViewData{}, err
		}
		user = models.NewRecord(users)
		user.Id = "nil"
	}

	if record, err := daoxtns.LoadDeckById(dao, c.PathParam("id")); err != nil {
		return ViewData{}, err
	} else {
		return ViewData{
			Title:       record.GetString("name"),
			Description: record.GetString("code") + "\n\n" + record.GetString("description"),
			Url:         c.Path(),
			Image:       record.GetString("image"),
			ImageAlt:    record.GetString("name"),
			Domain:      settings.Meta.AppUrl,
			Content: DeckPageData{
				Deck:        record,
				CurrentUser: user,
			},
		}, nil
	}
}

func LoadEditDeck(c echo.Context, dao *daos.Dao, settings *settings.Settings) (ViewData, error) {
	allSpellslingers, err := daoxtns.LoadSpellslingers(dao)
	if err != nil {
		return ViewData{}, err
	}
	colours, err := daoxtns.FindAllRecords(dao, "colours", nil)
	if err != nil {
		return ViewData{}, err
	}

	if record, err := daoxtns.LoadDeckById(dao, c.PathParam("id")); err != nil {
		return ViewData{}, err
	} else {
		return ViewData{
			Title:       record.GetString("name"),
			Description: record.GetString("code") + "\n\n" + record.GetString("description"),
			Url:         c.Path(),
			Image:       record.GetString("image"),
			ImageAlt:    record.GetString("name"),
			Domain:      settings.Meta.AppUrl,
			Content: EditDeckPageData{
				Deck:          record,
				Spellslingers: allSpellslingers,
				Colours:       colours,
			},
		}, nil
	}
}

func LoadEmptyDeck(c echo.Context, dao *daos.Dao, settings *settings.Settings) (ViewData, error) {
	user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
	if user == nil {
		return ViewData{}, errors.New("Missing auth record context.")
	}

	decks, err := dao.FindCollectionByNameOrId("decks")
	if err != nil {
		return ViewData{}, err
	}
	cards, err := dao.FindCollectionByNameOrId("cards")
	if err != nil {
		return ViewData{}, err
	}
	sets, err := dao.FindCollectionByNameOrId("sets")
	if err != nil {
		return ViewData{}, err
	}
	spellslingers, err := dao.FindCollectionByNameOrId("spellslingers")
	if err != nil {
		return ViewData{}, err
	}
	allSpellslingers, err := daoxtns.LoadSpellslingers(dao)
	if err != nil {
		return ViewData{}, err
	}
	colours, err := daoxtns.FindAllRecords(dao, "colours", nil)
	if err != nil {
		return ViewData{}, err
	}

	record := models.NewRecord(decks)
	land := models.NewRecord(cards)
	set := models.NewRecord(sets)
	spellslinger := models.NewRecord(spellslingers)

	set.Set("name", "default")
	land.Set("name", "default")
	land.Set("expand", map[string]any{
		"set": set,
	})
	spellslinger.Set("name", "default")

	expand := map[string]any{
		"cards":        []*models.Record{},
		"land":         land,
		"spellslinger": spellslinger,
	}

	record.Set("image", "/images/cards/tiles/default.jpeg")
	record.Set("cards", []*models.Record{})
	record.Set("name", "New Deck")
	record.Set("owner", user.Id)
	record.Set("is_public", false)
	record.Set("description", "")
	record.Set("code", "")
	record.Set("card_details", "")
	record.Set("expand", expand)

	return ViewData{
		Title:       "New Deck",
		Description: "",
		Url:         c.Path(),
		Image:       "",
		ImageAlt:    "",
		Domain:      settings.Meta.AppUrl,
		Content: EditDeckPageData{
			Deck:          record,
			Spellslingers: allSpellslingers,
			Colours:       colours,
		},
	}, nil
}

func LoadCard(c echo.Context, dao *daos.Dao, settings *settings.Settings) (ViewData, error) {
	if record, err := daoxtns.LoadCardById(dao, c.PathParam("id")); err != nil {
		return ViewData{}, err
	} else {
		expand := record.Expand()
		getExpand := func(name string) string {
			return expand[name].(*models.Record).GetString("name")
		}
		getExpands := func(name string) []string {
			records := expand[name].([]*models.Record)
			names := make([]string, len(records))
			for i, v := range records {
				names[i] = v.GetString("name")
			}
			return names
		}
		setName := getExpand("set")
		cardType := getExpand("type")
		ability := record.GetString("ability")

		descriptionLines := make([]string, 0)

		if cardType != "Land" && cardType != "Skill" {
			descriptionLines = append(descriptionLines, "Cost: "+record.GetString("cost"))
		}

		colours := getExpands("colour")
		descriptionLines = append(descriptionLines, strings.Join(colours, " "))

		if ability != "" {
			descriptionLines = append(descriptionLines, "\n"+ability+"\n")
		}

		descriptionLines = append(descriptionLines, cardType)

		if cardType == "Creature" {
			subTypes := getExpands("subtype")
			descriptionLines = append(descriptionLines, strings.Join(subTypes, " "))
			descriptionLines = append(descriptionLines, record.GetString("power")+"/"+record.GetString("health"))
		}

		if cardType == "Artifact" {
			descriptionLines = append(descriptionLines, "Charges: "+record.GetString("charges"))
		}

		if cardType == "Land" {
			descriptionLines = append(descriptionLines, "Chance: "+record.GetString("chance"))
		}

		descriptionLines = append(descriptionLines, setName+" "+getExpand("rarity"))
		joinedDescription := strings.Join(descriptionLines, "\n")

		description := ""
		if getExpand("rarity") == "Token" {
			description = joinedDescription
		}

		imageFormat := "text"
		switch c.QueryParam("image") {
		case "foil_full_art", "foil_text", "full_art":
			imageFormat = c.QueryParam("image")
			break
		}

		return ViewData{
			Title:       record.GetString("name"),
			Description: description,
			Url:         c.Path(),
			Image:       fmt.Sprintf("/images/cards/%s/%s.webp", imageFormat, record.GetString("name")),
			ImageAlt:    joinedDescription,
			Domain:      settings.Meta.AppUrl,
			Content:     record,
		}, nil
	}
}

func LoadSpellslinger(c echo.Context, dao *daos.Dao, settings *settings.Settings) (ViewData, error) {
	if record, err := daoxtns.LoadSpellslingerById(dao, c.PathParam("id")); err != nil {
		return ViewData{}, err
	} else {
		return ViewData{
			Title:       record.GetString("name"),
			Description: record.GetString("name"),
			Url:         c.Path(),
			Image:       fmt.Sprintf("/images/spellslingers/%s.webp", record.GetString("name")),
			ImageAlt:    record.GetString("name"),
			Domain:      settings.Meta.AppUrl,
			Content:     record,
		}, nil
	}
}

func LoadSpellslingers(c echo.Context, dao *daos.Dao, settings *settings.Settings) (ViewData, error) {
	if records, err := daoxtns.LoadSpellslingers(dao); err != nil {
		return ViewData{}, err
	} else {
		return ViewData{
			Title:       "Spellslingerer",
			Description: "A card database and deckbuilder companion to MtG Spellslingers",
			Url:         c.Path(),
			Image:       "/images/Spellslingerer.webp",
			ImageAlt:    "Spellslingerer.com logo",
			Domain:      settings.Meta.AppUrl,
			Content:     records,
		}, nil
	}
}

func LoadUser(c echo.Context, dao *daos.Dao, settings *settings.Settings, filter func(q *dbx.SelectQuery) error) (ViewData, error) {
	user, err := daoxtns.LoadUserById(dao, c.PathParam("userid"))
	if err != nil {
		return ViewData{}, err
	}

	return ViewData{
		Title:       fmt.Sprintf("%s's Decks", user.GetString("display_name")),
		Description: "A card database and deckbuilder companion to MtG Spellslingers",
		Url:         c.Path(),
		Image:       "/images/Spellslingerer.webp",
		ImageAlt:    "Spellslingerer.com logo",
		Domain:      settings.Meta.AppUrl,
		Content: DeckListPageData{
			Title:  fmt.Sprintf("%s's Decks", user.GetString("display_name")),
			Filter: template.JS(fmt.Sprintf("this.defaultFilter = `owner = '%s'`", c.PathParam("userid"))),
		},
	}, nil
}

func LoadCurrentUser(c echo.Context, dao *daos.Dao, settings *settings.Settings, filter func(q *dbx.SelectQuery) error) (ViewData, error) {
	user := c.Get(apis.ContextAuthRecordKey).(*models.Record)

	return ViewData{
		Title:       "Spellslingerer",
		Description: "A card database and deckbuilder companion to MtG Spellslingers",
		Url:         c.Path(),
		Image:       "/images/Spellslingerer.webp",
		ImageAlt:    "Spellslingerer.com logo",
		Domain:      settings.Meta.AppUrl,
		Content:     user,
	}, nil
}
