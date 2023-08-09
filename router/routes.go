package router

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"spellslingerer.com/m/v2/authxtns"
	m "spellslingerer.com/m/v2/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func BindRoutes(app core.App, router *echo.Echo) {
	dao := app.Dao()
	settings := app.Settings()

	router.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, c.Path(), DefaultViewData(c, settings))
	})

	router.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, c.Path(), DefaultViewData(c, settings))
	}, RequireGuestOnly())

	router.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, c.Path(), DefaultViewData(c, settings))
	}, RequireGuestOnly())

	router.GET("/forgot-password", func(c echo.Context) error {
		return c.Render(http.StatusOK, c.Path(), DefaultViewData(c, settings))
	}, RequireGuestOnly())

	router.GET("/profile", func(c echo.Context) error {
		if vd, err := LoadCurrentUser(c, dao, settings, nil); err != nil {
			return err
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	}, RequireRecordAuth())

	router.GET("/logout", func(c echo.Context) error {
		return authxtns.SetInvalidCookie(c, settings)
	}, RequireRecordAuth())

	router.GET("/cards", func(c echo.Context) error {
		return c.Render(http.StatusOK, c.Path(), DefaultViewData(c, settings))
	})

	router.GET("/cards/:id", func(c echo.Context) error {
		if vd, err := LoadCard(c, dao, settings); err != nil {
			return err
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	})

	router.GET("/spellslingers", func(c echo.Context) error {
		if vd, err := LoadSpellslingers(c, dao, settings); err != nil {
			return err
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	})

	router.GET("/spellslingers/:id", func(c echo.Context) error {
		if vd, err := LoadSpellslinger(c, dao, settings); err != nil {
			return err
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	})

	router.GET("/decks", func(c echo.Context) error {
		vd := DefaultViewData(c, settings)
		vd.Content = DeckListPageData{
			Title:  "Deck Search",
			Filter: "",
		}
		return c.Render(http.StatusOK, c.Path(), vd)
	})

	router.GET("/decks/:id", func(c echo.Context) error {
		if vd, err := LoadDeck(c, dao, settings); err != nil {
			return err
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	})

	router.GET("/decks/new", func(c echo.Context) error {
		if vd, err := LoadEmptyDeck(c, dao, settings); err != nil {
			return err
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	})

	router.POST("/decks/import", func(c echo.Context) error {
		code := c.FormValue("code")
		currentUser := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		deckId, err := m.ImportDeckByCode(app, code, currentUser.Id)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/decks/%s/edit", deckId))
	}, RequireRecordAuth())

	router.GET("/my-decks", func(c echo.Context) error {
		vd := DefaultViewData(c, settings)
		vd.Content = DeckListPageData{
			Title:  "My Decks",
			Filter: "this.defaultFilter = `owner = '${client.authStore.model.id}'`",
		}
		return c.Render(http.StatusOK, c.Path(), vd)
	}, RequireRecordAuth())

	router.GET("/decks/:id/edit", func(c echo.Context) error {
		currentUser := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		vd, err := LoadEditDeck(c, dao, settings)
		if err != nil {
			return err
		} else if currentUser.Id != vd.Content.(EditDeckPageData).Deck.GetString("owner") {
			return c.Render(http.StatusForbidden, strconv.Itoa(http.StatusForbidden), DefaultViewData(c, settings))
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	}, RequireRecordAuth())

	router.GET("/users/:userid", func(c echo.Context) error {
		if vd, err := LoadUser(c, dao, settings, nil); err != nil {
			return err
		} else {
			return c.Render(http.StatusOK, c.Path(), vd)
		}
	})

	router.GET("/integrations", func(c echo.Context) error {
		return c.Render(http.StatusOK, c.Path(), DefaultViewData(c, settings))
	})

	router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
}
