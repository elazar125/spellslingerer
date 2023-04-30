package router

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
)

func RequireGuestOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := echo.ErrBadRequest

			record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if record != nil {
				return err
			}

			admin, _ := c.Get(apis.ContextAdminKey).(*models.Admin)
			if admin != nil {
				return err
			}

			return next(c)
		}
	}
}

func RequireRecordAuth(optCollectionNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if record == nil {
				return echo.ErrUnauthorized
			}

			if len(optCollectionNames) > 0 && !list.ExistInSlice(record.Collection().Name, optCollectionNames) {
				return echo.ErrForbidden
			}

			return next(c)
		}
	}
}
