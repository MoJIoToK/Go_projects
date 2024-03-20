package endpoint

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	DaysLeft() int64
}

type Endpoint struct {
	s Service
}

// Constructor for Endpoint.
func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

// Status is the method which calls a function to count the number of days until 1/1/2025.
func (e *Endpoint) Status(ctx echo.Context) error {

	d := e.s.DaysLeft()

	s := fmt.Sprintf("Days left: %d", d)

	err := ctx.String(http.StatusOK, s)
	if err != nil {
		return err
	}

	return nil
}
