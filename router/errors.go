package router

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func GetErrorHandler(app core.App) echo.HTTPErrorHandler {
	return func(c echo.Context, err error) {
		if c.Response().Committed {
			return
		}

		var apiErr *apis.ApiError

		switch v := err.(type) {
		case *echo.HTTPError:
			msg := fmt.Sprintf("%v", v.Message)
			if err := c.Render(v.Code, strconv.Itoa(v.Code), DefaultViewData(c, app.Settings())); err != nil {
				log.Print(err)
			}
			apiErr = apis.NewApiError(v.Code, msg, v)
		case *apis.ApiError:
			apiErr = v
			if err := c.JSON(apiErr.Code, apiErr); err != nil {
				log.Print(err)
			}
		default:
			apiErr = apis.NewBadRequestError("", err)
			if err := c.JSON(apiErr.Code, apiErr); err != nil {
				log.Print(err)
			}
		}

		event := new(core.ApiErrorEvent)
		event.HttpContext = c
		event.Error = apiErr

		app.OnBeforeApiError().Trigger(event, func(e *core.ApiErrorEvent) error {
			// @see https://github.com/labstack/echo/issues/608
			if e.HttpContext.Request().Method == http.MethodHead {
				return e.HttpContext.NoContent(apiErr.Code)
			}

			return nil
		})

		app.OnAfterApiError().Trigger(event)
	}
}
