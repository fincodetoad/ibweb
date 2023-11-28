package ibweb

import (
	"encoding/json"
	"net/http"
)

const (
	positionByConIDPath = "v1/api/portfolio/{accountId}/position/{conid}"
)

type Position struct {
	AcctID            string   `json:"acctId"`
	Conid             int      `json:"conid"`
	ContractDesc      string   `json:"contractDesc"`
	AssetClass        string   `json:"assetClass"`
	Position          float64  `json:"position"`
	MktPrice          float64  `json:"mktPrice"`
	MktValue          float64  `json:"mktValue"`
	Currency          string   `json:"currency"`
	AvgCost           float64  `json:"avgCost"`
	AvgPrice          float64  `json:"avgPrice"`
	RealizedPnl       float64  `json:"realizedPnl"`
	UnrealizedPnl     float64  `json:"unrealizedPnl"`
	Exchs             string   `json:"exchs"`
	Expiry            string   `json:"expiry"`
	PutOrCall         string   `json:"putOrCall"`
	Multiplier        float64  `json:"multiplier"`
	Strike            string   `json:"strike"`
	ExerciseStyle     string   `json:"exerciseStyle"`
	UndConid          int      `json:"undConid"`
	ConExchMap        []string `json:"conExchMap"`
	BaseMktValue      float64  `json:"baseMktValue"`
	BaseMktPrice      float64  `json:"baseMktPrice"`
	BaseAvgCost       float64  `json:"baseAvgCost"`
	BaseAvgPrice      float64  `json:"baseAvgPrice"`
	BaseRealizedPnl   float64  `json:"baseRealizedPnl"`
	BaseUnrealizedPnl float64  `json:"baseUnrealizedPnl"`
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
}

/*
PositionsByContractID - Gets positions by contract ID
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1%7BaccountId%7D~1ledger/get
*/
func (c *client) PositionByContractID(accountID, conID string) ([]Position, error) {
	resp, err := c.get(substituteParam(positionByConIDPath,
		param{
			key:   "conid",
			value: conID,
		},
		param{
			key:   "accountId",
			value: accountID,
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

	var positionByConID []Position
	if err := json.Unmarshal(v, &positionByConID); err != nil {
		return nil, err
	}

	return positionByConID, nil
}
