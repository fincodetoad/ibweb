package ibweb

import (
	"encoding/json"
	"net/http"
)

const (
	marketDataHistory = "v1/api/iserver/marketdata/history"
)

type MarketDataHistory struct {
	ServerID           string `json:"serverId"`
	Symbol             string `json:"symbol"`
	Text               string `json:"text"`
	PriceFactor        int    `json:"priceFactor"`
	StartTime          string `json:"startTime"`
	High               string `json:"high"`
	Low                string `json:"low"`
	TimePeriod         string `json:"timePeriod"`
	BarLength          int    `json:"barLength"`
	MdAvailability     string `json:"mdAvailability"`
	MktDataDelay       int    `json:"mktDataDelay"`
	OutsideRth         bool   `json:"outsideRth"`
	TradingDayDuration int    `json:"tradingDayDuration"`
	VolumeFactor       int    `json:"volumeFactor"`
	PriceDisplayRule   int    `json:"priceDisplayRule"`
	PriceDisplayValue  string `json:"priceDisplayValue"`
	NegativeCapable    bool   `json:"negativeCapable"`
	MessageVersion     int    `json:"messageVersion"`
	Data               []struct {
		O float64 `json:"o"`
		C float64 `json:"c"`
		H float64 `json:"h"`
		L float64 `json:"l"`
		V int     `json:"v"`
		T int64   `json:"t"`
	} `json:"data"`
	Points     int `json:"points"`
	TravelTime int `json:"travelTime"`
}

type MarketDataHistoryInput struct {
	ConID      string
	Exchange   string
	Period     string
	Bar        string
	OutsideRth bool
}

func (m MarketDataHistoryInput) toQuery() []query {
	queries := []query{
		{
			key:   "conid",
			value: m.ConID,
		},
	}

	if m.Exchange != "" {
		queries = append(queries, query{
			key:   "exchange",
			value: m.Exchange,
		})
	}

	if m.Period != "" {
		queries = append(queries, query{
			key:   "period",
			value: m.Period,
		})
	}

	if m.Bar != "" {
		queries = append(queries, query{
			key:   "bar",
			value: m.Bar,
		})
	}

	if m.OutsideRth {
		queries = append(queries, query{
			key:   "outsideRth",
			value: "true",
		})
	}

	return queries
}

/*
MarketDataHistory - Gets the market data snapshot of a contract
Link: https://www.interactivebrokers.com/api/doc.html#tag/Market-Data/paths/~1iserver~1marketdata~1history/get
*/
func (c *client) MarketDataHistory(input MarketDataHistoryInput) (*MarketDataHistory, error) {
	resp, err := c.get(marketDataHistory, input.toQuery()...)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, StatusCodeError{StatusCode: resp.StatusCode, Err: NewIBError(resp)}
	}

	defer resp.Body.Close()
	v, err := readAllFn(resp.Body)
	if err != nil {
		return nil, err
	}

	var marketDataHistory MarketDataHistory
	if err := json.Unmarshal(v, &marketDataHistory); err != nil {
		return nil, err
	}

	return &marketDataHistory, nil
}
