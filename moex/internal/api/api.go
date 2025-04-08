package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"moex-app/internal/model"
	"net/http"
)

func GetMarketData(ticker string) (*model.ISSResponse, error) {

	log := slog.Default()
	log.Info("GetMarketData")

	var result model.ISSResponse

	url := fmt.Sprintf(
		"https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities.json?securities=%s",
		ticker,
	)

	log.Info("Request URL: " + url)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(resp.Status)
		return nil, fmt.Errorf("bad status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("read body error: %v", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}

	return &result, nil
}
