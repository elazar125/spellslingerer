package authxtns

import (
	"net/http"
	"time"

	"spellslingerer.com/m/v2/config"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models/settings"
)

func LoadCookieContext(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, _ := c.Cookie(config.COOKIE_NAME)

			if cookie != nil {
				c.Request().Header.Add("Authorization", cookie.Value)
			}

			return next(c)
		}
	}
}

func getDefaultCookie(settings *settings.Settings) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = config.COOKIE_NAME
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Secure = true
	cookie.Domain = settings.Meta.AppUrl
	cookie.Path = "/"
	return cookie
}

func SetValidCookieHandler(settings *settings.Settings) func(*core.RecordAuthEvent) error {
	return func(e *core.RecordAuthEvent) error {
		cookie := getDefaultCookie(settings)
		cookie.Value = e.Token
		cookie.Expires = time.Now().Add(time.Duration(settings.RecordAuthToken.Duration) * time.Second)
		e.HttpContext.SetCookie(cookie)
		return nil
	}
}

func SetInvalidCookie(c echo.Context, settings *settings.Settings) error {
	cookie := getDefaultCookie(settings)
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1)
	c.SetCookie(cookie)
	return nil
}
