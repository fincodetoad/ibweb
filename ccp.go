package ibweb

import (
	"encoding/json"
	"net/http"
)

const (
	positionsByConIDPath = "v1/api/portfolio/positions/{conid}"
)

type PositionsByConID struct {
	ACCTID []struct {
		AcctID            string   `json:"acctId"`
		Conid             int      `json:"conid"`
		ContractDesc      string   `json:"contractDesc"`
		AssetClass        string   `json:"assetClass"`
		Position          int      `json:"position"`
		MktPrice          int      `json:"mktPrice"`
		MktValue          int      `json:"mktValue"`
		Currency          string   `json:"currency"`
		AvgCost           int      `json:"avgCost"`
		AvgPrice          int      `json:"avgPrice"`
		RealizedPnl       int      `json:"realizedPnl"`
		UnrealizedPnl     int      `json:"unrealizedPnl"`
		Exchs             string   `json:"exchs"`
		Expiry            string   `json:"expiry"`
		PutOrCall         string   `json:"putOrCall"`
		Multiplier        int      `json:"multiplier"`
		Strike            int      `json:"strike"`
		ExerciseStyle     string   `json:"exerciseStyle"`
		UndConid          int      `json:"undConid"`
		ConExchMap        []string `json:"conExchMap"`
		BaseMktValue      int      `json:"baseMktValue"`
		BaseMktPrice      int      `json:"baseMktPrice"`
		BaseAvgCost       int      `json:"baseAvgCost"`
		BaseAvgPrice      int      `json:"baseAvgPrice"`
		BaseRealizedPnl   int      `json:"baseRealizedPnl"`
		BaseUnrealizedPnl int      `json:"baseUnrealizedPnl"`
		Name              string   `json:"name"`
		LastTradingDay    string   `json:"lastTradingDay"`
		Group             string   `json:"group"`
		Sector            string   `json:"sector"`
		SectorGroup       string   `json:"sectorGroup"`
		Ticker            string   `json:"ticker"`
		UndComp           string   `json:"undComp"`
		UndSym            string   `json:"undSym"`
		FullName          string   `json:"fullName"`
		PageSize          int      `json:"pageSize"`
		Model             string   `json:"model"`
	} `json:"ACCTID"`
}

/*
PositionsByContractID - Gets positions by contract ID
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1%7BaccountId%7D~1ledger/get
*/
func (c *client) PositionsByContractID(conID string) (*PositionsByConID, error) {
	resp, err := c.get(substituteParam(positionsByConIDPath,
		param{
			key:   "conID",
			value: conID,
		},
	))
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

	var positionByConID PositionsByConID
	if err := json.Unmarshal(v, &positionByConID); err != nil {
		return nil, err
	}

	return &positionByConID, nil
}
