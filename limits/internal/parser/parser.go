// Package parser reads raw input (e.g., from file or HTTP body) and extracts DEPO records.
// It supports the format: DEPO: FIRM_ID = #; SECCODE = #; CLIENT_CODE = #; ...
// Only DEPO lines are processed; MONEY lines are ignored.
package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"limits-app/internal/models"
	"strconv"
	"strings"
)

var requiredDepoLimit = []string{"CLIENT_CODE", "SECCODE", "LIMIT_KIND", "OPEN_LIMIT", "OPEN_BALANCE"}

// Parse reads an io.Reader line by line and extracts valid DEPO records.
// Returns a slice of DepoLimit and any parsing errors (non-fatal).
// Lines not starting with "DEPO:" are skipped.
func Parse(reader io.Reader) ([]models.DepoLimit, []error) {
	var depoLimits []models.DepoLimit
	var errs []error

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "DEPO:") {
			var depoLimit models.DepoLimit
			var lineErr []error
			depoLimit, lineErr = parseDepoLine(line)
			if lineErr != nil {
				errs = append(errs, lineErr...)
				fmt.Println(lineErr)
			}
			depoLimits = append(depoLimits, depoLimit)
		}
	}

	return depoLimits, errs
}

// parseDepoLine parses a single DEPO line into a DepoLimit struct.
// Expected format: "DEPO: KEY = VALUE; ..."
// Returns the parsed limit and any field-level errors.
func parseDepoLine(line string) (models.DepoLimit, []error) {
	var errs []error
	var err error
	rawLine := strings.TrimPrefix(line, "DEPO: ")
	rawLine = strings.TrimSpace(rawLine)

	depoPairs := strings.Split(rawLine, ";")
	resMap := make(map[string]string)

	for _, pair := range depoPairs {
		key, value, found := strings.Cut(pair, "=")
		if found {
			resMap[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}

	if !hasAllKeys(resMap) {
		errs = append(errs, errors.New("Depo limit does not contain all depo limits"))
	}

	depoLimit := models.DepoLimit{}
	depoLimit.ClientCode = resMap["CLIENT_CODE"]
	depoLimit.SecCode = resMap["SECCODE"]
	depoLimit.LimitKind, err = strconv.Atoi(resMap["LIMIT_KIND"])
	if err != nil {
		errs = append(errs, err)
	}
	depoLimit.OpenLimit, err = strconv.ParseFloat(resMap["OPEN_LIMIT"], 64)
	if err != nil {
		errs = append(errs, err)
	}
	depoLimit.OpenBalance, err = strconv.ParseFloat(resMap["OPEN_BALANCE"], 64)
	if err != nil {
		errs = append(errs, err)
	}

	return depoLimit, errs
}

// hasAllKeys checks that all required keys are present in the parsed map.
func hasAllKeys(resMap map[string]string) bool {
	for _, key := range requiredDepoLimit {
		if _, found := resMap[key]; !found {
			return false
		}
	}
	return true
}
