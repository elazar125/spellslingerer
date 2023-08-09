package authxtns

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

const COOKIE_NAME string = "pb_auth"

func LoadCookieContext(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(COOKIE_NAME)
			if err != nil || cookie.Value == "" {
				return next(c)
			}

			claims, _ := security.ParseUnverifiedJWT(cookie.Value)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
				admin, err := app.Dao().FindAdminByToken(
					cookie.Value,
					app.Settings().AdminAuthToken.Secret,
				)
				if err == nil && admin != nil {
					c.Set(apis.ContextAdminKey, admin)
				}
			case tokens.TypeAuthRecord:
				record, err := app.Dao().FindAuthRecordByToken(
					cookie.Value,
					app.Settings().RecordAuthToken.Secret,
				)
				if err == nil && record != nil {
					c.Set(apis.ContextAuthRecordKey, record)
				}
			}

			return next(c)
		}
	}
}

func getDefaultCookie(settings *settings.Settings) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = COOKIE_NAME
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
