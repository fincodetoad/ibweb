package ibweb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	searchContractsPath = "v1/api/iserver/secdef/search"
	searchStrikesPath   = "v1/api/iserver/secdef/strikes"
	secDefInfoPath      = "v1/api/iserver/secdef/info"
)

// SecType - Security type of contract
type SecType string

const (
	Options SecType = "OPT"
	Stock   SecType = "STK"
	War     SecType = "WAR"
)

// Right - Options right
type Right string

const (
	Call Right = "C"
	Put  Right = "P"
)

/*
SearchContractsInput -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1search/post
*/
type SearchContractsInput struct {
	Symbol  string  `json:"symbol"`
	Name    bool    `json:"name"`
	SecType SecType `json:"sectype"`
}

/*
Contract -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1search/post
*/
type Contract struct {
	Conid         string `json:"conid"`
	CompanyHeader string `json:"companyHeader"`
	CompanyName   string `json:"companyName"`
	Symbol        string `json:"symbol"`
	Description   string `json:"description"`
	Restricted    string `json:"restricted"`
	Fop           string `json:"fop"`
	Opt           string `json:"opt"`
	War           string `json:"war"`
	Sections      []struct {
		SecType    SecType `json:"secType"`
		Months     string  `json:"months"`
		Symbol     string  `json:"symbol"`
		Exchange   string  `json:"exchange"`
		LegSecType string  `json:"legSecType"`
	} `json:"sections"`
}

/*
SearchStrikesInput -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1strikes/get
*/
type SearchStrikesInput struct {
	ConID    string
	SecType  SecType
	Month    string
	Exchange string
}

func (s SearchStrikesInput) toQuery() []query {
	queries := []query{
		{
			key:   "conid",
			value: s.ConID,
		},
		{
			key:   "sectype",
			value: string(s.SecType),
		},
		{
			key:   "month",
			value: s.Month,
		},
	}

	if s.Exchange != "" {
		queries = append(queries, query{
			key:   "exchange",
			value: s.Exchange,
		})
	}

	return queries
}

/*
SearchStrikes -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1strikes/get
*/
type SearchStrikes struct {
	Call []float64 `json:"call"`
	Put  []float64 `json:"put"`
}

/*
SecurityDefinitionInfoInput -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1info/get
*/
type SecurityDefinitionInfoInput struct {
	ConID    string
	SecType  SecType
	Month    string
	Exchange string
	Strike   float64
	Right    string
}

func (s SecurityDefinitionInfoInput) toQuery() []query {
	queries := []query{
		{
			key:   "conid",
			value: s.ConID,
		},
		{
			key:   "sectype",
			value: string(s.SecType),
		},
	}

	if s.Month != "" {
		queries = append(queries, query{
			key:   "month",
			value: s.Month,
		})
	}

	if s.Exchange != "" {
		queries = append(queries, query{
			key:   "exchange",
			value: s.Exchange,
		})
	}

	if s.Strike != 0 {
		queries = append(queries, query{
			key:   "strike",
			value: fmt.Sprintf("%f", s.Strike),
		})
	}

	if s.Right != "" {
		queries = append(queries, query{
			key:   "right",
			value: s.Right,
		})
	}

	return queries
}

/*
SecurityDefinitionInfo -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1info/get
*/
type SecurityDefinitionInfo struct {
	Conid           int     `json:"conid"`
	Symbol          string  `json:"symbol"`
	SecType         string  `json:"secType"`
	Exchange        string  `json:"exchange"`
	ListingExchange string  `json:"listingExchange"`
	Right           string  `json:"right"`
	Strike          float64 `json:"strike"`
	Currency        string  `json:"currency"`
	Cusip           string  `json:"cusip"`
	Coupon          string  `json:"coupon"`
	Desc1           string  `json:"desc1"`
	Desc2           string  `json:"desc2"`
	MaturityDate    string  `json:"maturityDate"`
	Multiplier      string  `json:"multiplier"`
	TradingClass    string  `json:"tradingClass"`
	ValidExchanges  string  `json:"validExchanges"`
}

/*
SearchContracts - Searches for a contract
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1search/post
*/
func (c *client) SearchContracts(input SearchContractsInput) ([]Contract, error) {
	resp, err := c.post(searchContractsPath, &input)
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

	var contracts []Contract
	if err := json.Unmarshal(v, &contracts); err != nil {
		return nil, err
	}

	return contracts, nil
}

/*
SearchStrikes - Searches for the strike prices of a contract
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1strikes/get
*/
func (c *client) SearchStrikes(input SearchStrikesInput) (*SearchStrikes, error) {
	resp, err := c.get(searchStrikesPath, input.toQuery()...)
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

	var searchStrikes SearchStrikes
	if err := json.Unmarshal(v, &searchStrikes); err != nil {
		return nil, err
	}

	return &searchStrikes, nil
}

/*
SecurityDefinitionInfo - Get a security definition informaiton
Link: https://www.interactivebrokers.com/api/doc.html#tag/Contract/paths/~1iserver~1secdef~1info/get
*/
func (c *client) SecurityDefinitionInfo(input SecurityDefinitionInfoInput) ([]SecurityDefinitionInfo, error) {
	resp, err := c.get(secDefInfoPath, input.toQuery()...)
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

	var secDefInfo []SecurityDefinitionInfo
	if err := json.Unmarshal(v, &secDefInfo); err != nil {
		return nil, err
	}

	return secDefInfo, nil
}
