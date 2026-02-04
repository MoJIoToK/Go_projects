// Package normalize implements business logic for filtering and validating DEPO limits
// according to the task requirements.
package normalize

import (
	"fmt"
	"limits-app/internal/models"
)

// Input contains the raw depo limits received from the parser.
type Input struct {
	DepoLimits []models.DepoLimit
}

// Output contains normalized client db and non-fatal validation errors.
type Output struct {
	ClientLimits []models.ClientLimits
	Errors       []error
}

// Normalizer defines the interface for limit normalization.
type Normalizer interface {
	Normalize(input Input) (Output, error)
}

type normalizer struct {
}

// NewNormalizer is a constructor for the normalizer.
func NewNormalizer() Normalizer {
	return &normalizer{}
}

// Normalize processes input raw limits according to business rules:
// - Filters out zero-limit positions,
// - Validates presence of all 4 limits for non-zero positions.
// - Adds fake position if client becomes empty,
func (n *normalizer) Normalize(input Input) (Output, error) {
	clientMap := groupByClientAndSecID(input.DepoLimits)

	var results []models.ClientLimits
	var errs []error

	for client, positions := range clientMap {
		validPositions := filterZeroPositions(positions)

		errorsValidation := validateAllLimits(client, validPositions)
		errs = append(errs, errorsValidation...)

		if len(validPositions) == 0 && len(positions) > 0 {
			fakePosition := createFakePosition(client, positions)
			validPositions = append(validPositions, fakePosition)
		}
		results = append(results, models.ClientLimits{
			ClientCode: client,
			Positions:  validPositions,
		})

	}
	return Output{ClientLimits: results, Errors: errs}, nil
}

// createFakePosition generates a fake position for clients with no valid positions.
func createFakePosition(client string, positions map[string][]models.DepoLimit) models.InstrumentPosition {
	var secID string
	for id := range positions {
		secID = id
		break
	}
	if secID == "" {
		secID = "X5"
	}

	return models.InstrumentPosition{
		SecId: secID,
		Limits: map[int]models.DepoLimit{
			0: {
				ClientCode:  client,
				SecCode:     secID,
				LimitKind:   0,
				OpenLimit:   0,
				OpenBalance: 0,
			},
		},
	}
}

// filterZeroPositions removes positions with OPEN_BALANCE = 0 where all 4 limits have OPEN_LIMIT = 0.
func filterZeroPositions(positions map[string][]models.DepoLimit) []models.InstrumentPosition {
	var results []models.InstrumentPosition

	for secId, limits := range positions {
		limitMap := make(map[int]models.DepoLimit)
		allZero := true
		isZeroBalance := false

		for _, limit := range limits {
			limitMap[limit.LimitKind] = limit
			if limit.OpenLimit != 0 {
				allZero = false
			}
			if limit.OpenBalance == 0 {
				isZeroBalance = true
			}
		}

		if allZero && len(limitMap) == 4 && isZeroBalance {
			continue
		}

		position := models.InstrumentPosition{
			SecId:       secId,
			OpenBalance: limitMap[0].OpenBalance,
			Limits:      limitMap,
		}
		for key, value := range limitMap {
			position.Limits[key] = value
		}
		results = append(results, position)
	}
	return results
}

// groupByClientAndSecID groups depo limits by client code and security code.
func groupByClientAndSecID(depoLimits []models.DepoLimit) map[string]map[string][]models.DepoLimit {

	clientMap := make(map[string]map[string][]models.DepoLimit)
	for _, dl := range depoLimits {
		if clientMap[dl.ClientCode] == nil {
			clientMap[dl.ClientCode] = make(map[string][]models.DepoLimit)
		}
		clientMap[dl.ClientCode][dl.SecCode] = append(clientMap[dl.ClientCode][dl.SecCode], dl)
	}

	return clientMap
}

// validateAllLimits checks that every position with non-zero OpenBalance
// has all 4 required limit kinds (0, 1, 2, 365).
// Missing limits are reported as non-fatal errors.
func validateAllLimits(client string, positions []models.InstrumentPosition) []error {
	var errs []error
	requiredKinds := []int{0, 1, 2, 365}

	for _, pos := range positions {
		isNonZeroBalance := false
		for _, limit := range pos.Limits {
			if limit.OpenBalance != 0 {
				isNonZeroBalance = true
				break
			}
		}

		if !isNonZeroBalance {
			continue
		}

		for _, kind := range requiredKinds {
			if _, exists := pos.Limits[kind]; !exists {
				errs = append(errs, fmt.Errorf(
					"client %s, sec %s: missing LIMIT_KIND=%d",
					client, pos.SecId, kind))
			}
		}
	}
	return errs
}
