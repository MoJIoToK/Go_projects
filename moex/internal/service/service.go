package service

import (
	"errors"
	"log/slog"
	"moex-app/internal/model"
)

func ExtractData(doc *model.ISSResponse) ([]map[string]interface{}, error) {

	log := slog.Default()

	log.Info("ExtractData")

	if doc == nil {
		log.Error("doc is nil")
		return nil, errors.New("doc is nil")
	}

	secidToSecname := make(map[string]string)
	for _, row := range doc.Securities.Data {
		secid := row[0].(string)
		secName := row[9].(string)
		secidToSecname[secid] = secName
	}

	columns := doc.Marketdata.Columns
	result := make([]map[string]interface{}, 0, len(columns)+1)

	for _, row := range doc.Marketdata.Data {
		if row == nil {
			log.Error("row is nil")
			continue
		}
		
		entry := make(map[string]interface{}, len(columns)+1)
		for i := 0; i < len(columns) && i < len(row); i++ {
			entry[columns[i]] = row[i]
		}

		if secName, ok := secidToSecname[row[0].(string)]; ok {
			entry["SECNAME"] = secName
		}

		result = append(result, entry)
	}

	return result, nil
}
