package mw

import (
	"github.com/labstack/echo/v4"
	"log"
)

func CheckRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		//Getting the header value
		val := ctx.Request().Header.Get("User-Role")

		if val == "admin" {
			log.Println("red button user detected")
		}

		err := next(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}
