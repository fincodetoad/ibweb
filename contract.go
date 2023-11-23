package ibweb

import (
	"encoding/json"
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
	Options = "OPT"
	Stock   = "STK"
	War     = "WAR"
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
