package ibweb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

const (
	portfolioAccountsPath  = "v1/api/portfolio/accounts"
	subAccountsPath        = "v1/api/portfolio/subaccounts"
	subAccountsLargePath   = "v1/api/portfolio/subaccounts2"
	accountInformationPath = "v1/api/portfolio/{accountId}/meta"
	accountSummaryPath     = "v1/api/portfolio/{accountId}/summary"
)

/*
PortfolioAccount
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1accounts/get
*/
type PortfolioAccount struct {
	ID             string `json:"id"`
	AccountID      string `json:"accountId"`
	AccountVan     string `json:"accountVan"`
	AccountTitle   string `json:"accountTitle"`
	DisplayName    string `json:"displayName"`
	AccountAlias   string `json:"accountAlias"`
	AccountStatus  int    `json:"accountStatus"`
	Currency       string `json:"currency"`
	Type           string `json:"type"`
	TradingType    string `json:"tradingType"`
	Faclient       bool   `json:"faclient"`
	ClearingStatus string `json:"clearingStatus"`
	Covestor       bool   `json:"covestor"`
	Parent         struct {
		Mmc         []string `json:"mmc"`
		AccountID   string   `json:"accountId"`
		IsMParent   bool     `json:"isMParent"`
		IsMChild    bool     `json:"isMChild"`
		IsMultiplex bool     `json:"isMultiplex"`
	} `json:"parent"`
	Desc string `json:"desc"`
}

/*
SubAccount
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1subaccounts/get
*/
type SubAccount PortfolioAccount

/*
SubAccountsLarge
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1subaccounts2/get
*/
type SubAccountsLarge struct {
	Metadata struct {
		Total    int `json:"total"`
		PageSize int `json:"pageSize"`
		PageNume int `json:"pageNume"`
	} `json:"metadata"`
	Subaccounts []SubAccount `json:"subaccounts"`
}

/*
AccountInformation
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1%7BaccountId%7D~1meta/get
*/
type AccountInformation struct {
	ID             string `json:"id"`
	AccountID      string `json:"accountId"`
	AccountVan     string `json:"accountVan"`
	AccountTitle   string `json:"accountTitle"`
	DisplayName    string `json:"displayName"`
	AccountAlias   string `json:"accountAlias"`
	AccountStatus  int    `json:"accountStatus"`
	Currency       string `json:"currency"`
	Type           string `json:"type"`
	TradingType    string `json:"tradingType"`
	Faclient       bool   `json:"faclient"`
	ClearingStatus string `json:"clearingStatus"`
	Covestor       bool   `json:"covestor"`
	Parent         struct {
		Mmc         []string `json:"mmc"`
		AccountID   string   `json:"accountId"`
		IsMParent   bool     `json:"isMParent"`
		IsMChild    bool     `json:"isMChild"`
		IsMultiplex bool     `json:"isMultiplex"`
	} `json:"parent"`
	Desc string `json:"desc"`
}

/*
AccountSummary
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1%7BaccountId%7D~1summary/get
*/
type AccountSummary struct {
	Accountready                    AccountSummaryInner `json:"accountready"`
	Accounttype                     AccountSummaryInner `json:"accounttype"`
	Accruedcash                     AccountSummaryInner `json:"accruedcash"`
	AccruedcashC                    AccountSummaryInner `json:"accruedcash-c"`
	AccruedcashF                    AccountSummaryInner `json:"accruedcash-f"`
	AccruedcashS                    AccountSummaryInner `json:"accruedcash-s"`
	Accrueddividend                 AccountSummaryInner `json:"accrueddividend"`
	AccrueddividendC                AccountSummaryInner `json:"accrueddividend-c"`
	AccrueddividendF                AccountSummaryInner `json:"accrueddividend-f"`
	AccrueddividendS                AccountSummaryInner `json:"accrueddividend-s"`
	Availablefunds                  AccountSummaryInner `json:"availablefunds"`
	AvailablefundsC                 AccountSummaryInner `json:"availablefunds-c"`
	AvailablefundsF                 AccountSummaryInner `json:"availablefunds-f"`
	AvailablefundsS                 AccountSummaryInner `json:"availablefunds-s"`
	Billable                        AccountSummaryInner `json:"billable"`
	BillableC                       AccountSummaryInner `json:"billable-c"`
	BillableF                       AccountSummaryInner `json:"billable-f"`
	BillableS                       AccountSummaryInner `json:"billable-s"`
	Buyingpower                     AccountSummaryInner `json:"buyingpower"`
	Cushion                         AccountSummaryInner `json:"cushion"`
	Daytradesremaining              AccountSummaryInner `json:"daytradesremaining"`
	Daytradesremainingt1            AccountSummaryInner `json:"daytradesremainingt+1"`
	Daytradesremainingt2            AccountSummaryInner `json:"daytradesremainingt+2"`
	Daytradesremainingt3            AccountSummaryInner `json:"daytradesremainingt+3"`
	Daytradesremainingt4            AccountSummaryInner `json:"daytradesremainingt+4"`
	Equitywithloanvalue             AccountSummaryInner `json:"equitywithloanvalue"`
	EquitywithloanvalueC            AccountSummaryInner `json:"equitywithloanvalue-c"`
	EquitywithloanvalueF            AccountSummaryInner `json:"equitywithloanvalue-f"`
	EquitywithloanvalueS            AccountSummaryInner `json:"equitywithloanvalue-s"`
	Excessliquidity                 AccountSummaryInner `json:"excessliquidity"`
	ExcessliquidityC                AccountSummaryInner `json:"excessliquidity-c"`
	ExcessliquidityF                AccountSummaryInner `json:"excessliquidity-f"`
	ExcessliquidityS                AccountSummaryInner `json:"excessliquidity-s"`
	Fullavailablefunds              AccountSummaryInner `json:"fullavailablefunds"`
	FullavailablefundsC             AccountSummaryInner `json:"fullavailablefunds-c"`
	FullavailablefundsF             AccountSummaryInner `json:"fullavailablefunds-f"`
	FullavailablefundsS             AccountSummaryInner `json:"fullavailablefunds-s"`
	Fullexcessliquidity             AccountSummaryInner `json:"fullexcessliquidity"`
	FullexcessliquidityC            AccountSummaryInner `json:"fullexcessliquidity-c"`
	FullexcessliquidityF            AccountSummaryInner `json:"fullexcessliquidity-f"`
	FullexcessliquidityS            AccountSummaryInner `json:"fullexcessliquidity-s"`
	Fullinitmarginreq               AccountSummaryInner `json:"fullinitmarginreq"`
	FullinitmarginreqC              AccountSummaryInner `json:"fullinitmarginreq-c"`
	FullinitmarginreqF              AccountSummaryInner `json:"fullinitmarginreq-f"`
	FullinitmarginreqS              AccountSummaryInner `json:"fullinitmarginreq-s"`
	Fullmaintmarginreq              AccountSummaryInner `json:"fullmaintmarginreq"`
	FullmaintmarginreqC             AccountSummaryInner `json:"fullmaintmarginreq-c"`
	FullmaintmarginreqF             AccountSummaryInner `json:"fullmaintmarginreq-f"`
	FullmaintmarginreqS             AccountSummaryInner `json:"fullmaintmarginreq-s"`
	Grosspositionvalue              AccountSummaryInner `json:"grosspositionvalue"`
	GrosspositionvalueC             AccountSummaryInner `json:"grosspositionvalue-c"`
	GrosspositionvalueF             AccountSummaryInner `json:"grosspositionvalue-f"`
	GrosspositionvalueS             AccountSummaryInner `json:"grosspositionvalue-s"`
	Guarantee                       AccountSummaryInner `json:"guarantee"`
	GuaranteeC                      AccountSummaryInner `json:"guarantee-c"`
	GuaranteeF                      AccountSummaryInner `json:"guarantee-f"`
	GuaranteeS                      AccountSummaryInner `json:"guarantee-s"`
	Highestseverity                 AccountSummaryInner `json:"highestseverity"`
	HighestseverityC                AccountSummaryInner `json:"highestseverity-c"`
	HighestseverityF                AccountSummaryInner `json:"highestseverity-f"`
	HighestseverityS                AccountSummaryInner `json:"highestseverity-s"`
	Indianstockhaircut              AccountSummaryInner `json:"indianstockhaircut"`
	IndianstockhaircutC             AccountSummaryInner `json:"indianstockhaircut-c"`
	IndianstockhaircutF             AccountSummaryInner `json:"indianstockhaircut-f"`
	IndianstockhaircutS             AccountSummaryInner `json:"indianstockhaircut-s"`
	Initmarginreq                   AccountSummaryInner `json:"initmarginreq"`
	InitmarginreqC                  AccountSummaryInner `json:"initmarginreq-c"`
	InitmarginreqF                  AccountSummaryInner `json:"initmarginreq-f"`
	InitmarginreqS                  AccountSummaryInner `json:"initmarginreq-s"`
	Leverage                        AccountSummaryInner `json:"leverage"`
	LeverageC                       AccountSummaryInner `json:"leverage-c"`
	LeverageF                       AccountSummaryInner `json:"leverage-f"`
	LeverageS                       AccountSummaryInner `json:"leverage-s"`
	Lookaheadavailablefunds         AccountSummaryInner `json:"lookaheadavailablefunds"`
	LookaheadavailablefundsC        AccountSummaryInner `json:"lookaheadavailablefunds-c"`
	LookaheadavailablefundsF        AccountSummaryInner `json:"lookaheadavailablefunds-f"`
	LookaheadavailablefundsS        AccountSummaryInner `json:"lookaheadavailablefunds-s"`
	Lookaheadexcessliquidity        AccountSummaryInner `json:"lookaheadexcessliquidity"`
	LookaheadexcessliquidityC       AccountSummaryInner `json:"lookaheadexcessliquidity-c"`
	LookaheadexcessliquidityF       AccountSummaryInner `json:"lookaheadexcessliquidity-f"`
	LookaheadexcessliquidityS       AccountSummaryInner `json:"lookaheadexcessliquidity-s"`
	Lookaheadinitmarginreq          AccountSummaryInner `json:"lookaheadinitmarginreq"`
	LookaheadinitmarginreqC         AccountSummaryInner `json:"lookaheadinitmarginreq-c"`
	LookaheadinitmarginreqF         AccountSummaryInner `json:"lookaheadinitmarginreq-f"`
	LookaheadinitmarginreqS         AccountSummaryInner `json:"lookaheadinitmarginreq-s"`
	Lookaheadmaintmarginreq         AccountSummaryInner `json:"lookaheadmaintmarginreq"`
	LookaheadmaintmarginreqC        AccountSummaryInner `json:"lookaheadmaintmarginreq-c"`
	LookaheadmaintmarginreqF        AccountSummaryInner `json:"lookaheadmaintmarginreq-f"`
	LookaheadmaintmarginreqS        AccountSummaryInner `json:"lookaheadmaintmarginreq-s"`
	Lookaheadnextchange             AccountSummaryInner `json:"lookaheadnextchange"`
	Maintmarginreq                  AccountSummaryInner `json:"maintmarginreq"`
	MaintmarginreqC                 AccountSummaryInner `json:"maintmarginreq-c"`
	MaintmarginreqF                 AccountSummaryInner `json:"maintmarginreq-f"`
	MaintmarginreqS                 AccountSummaryInner `json:"maintmarginreq-s"`
	Netliquidation                  AccountSummaryInner `json:"netliquidation"`
	NetliquidationC                 AccountSummaryInner `json:"netliquidation-c"`
	NetliquidationF                 AccountSummaryInner `json:"netliquidation-f"`
	NetliquidationS                 AccountSummaryInner `json:"netliquidation-s"`
	Netliquidationuncertainty       AccountSummaryInner `json:"netliquidationuncertainty"`
	Nlvandmargininreview            AccountSummaryInner `json:"nlvandmargininreview"`
	Pasharesvalue                   AccountSummaryInner `json:"pasharesvalue"`
	PasharesvalueC                  AccountSummaryInner `json:"pasharesvalue-c"`
	PasharesvalueF                  AccountSummaryInner `json:"pasharesvalue-f"`
	PasharesvalueS                  AccountSummaryInner `json:"pasharesvalue-s"`
	Postexpirationexcess            AccountSummaryInner `json:"postexpirationexcess"`
	PostexpirationexcessC           AccountSummaryInner `json:"postexpirationexcess-c"`
	PostexpirationexcessF           AccountSummaryInner `json:"postexpirationexcess-f"`
	PostexpirationexcessS           AccountSummaryInner `json:"postexpirationexcess-s"`
	Postexpirationmargin            AccountSummaryInner `json:"postexpirationmargin"`
	PostexpirationmarginC           AccountSummaryInner `json:"postexpirationmargin-c"`
	PostexpirationmarginF           AccountSummaryInner `json:"postexpirationmargin-f"`
	PostexpirationmarginS           AccountSummaryInner `json:"postexpirationmargin-s"`
	Previousdayequitywithloanvalue  AccountSummaryInner `json:"previousdayequitywithloanvalue"`
	PreviousdayequitywithloanvalueC AccountSummaryInner `json:"previousdayequitywithloanvalue-c"`
	PreviousdayequitywithloanvalueF AccountSummaryInner `json:"previousdayequitywithloanvalue-f"`
	PreviousdayequitywithloanvalueS AccountSummaryInner `json:"previousdayequitywithloanvalue-s"`
	SegmenttitleC                   AccountSummaryInner `json:"segmenttitle-c"`
	SegmenttitleF                   AccountSummaryInner `json:"segmenttitle-f"`
	SegmenttitleS                   AccountSummaryInner `json:"segmenttitle-s"`
	Totalcashvalue                  AccountSummaryInner `json:"totalcashvalue"`
	TotalcashvalueC                 AccountSummaryInner `json:"totalcashvalue-c"`
	TotalcashvalueF                 AccountSummaryInner `json:"totalcashvalue-f"`
	TotalcashvalueS                 AccountSummaryInner `json:"totalcashvalue-s"`
	Totaldebitcardpendingcharges    AccountSummaryInner `json:"totaldebitcardpendingcharges"`
	TotaldebitcardpendingchargesC   AccountSummaryInner `json:"totaldebitcardpendingcharges-c"`
	TotaldebitcardpendingchargesF   AccountSummaryInner `json:"totaldebitcardpendingcharges-f"`
	TotaldebitcardpendingchargesS   AccountSummaryInner `json:"totaldebitcardpendingcharges-s"`
	TradingtypeF                    AccountSummaryInner `json:"tradingtype-f"`
	TradingtypeS                    AccountSummaryInner `json:"tradingtype-s"`
}

type AccountSummaryInner struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	IsNull    bool    `json:"isNull"`
	Timestamp int     `json:"timestamp"`
	Value     float64 `json:"value"`
}

func (a *AccountSummaryInner) UnmarshalJSON(v []byte) error {
	type altStruct struct {
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		IsNull    bool    `json:"isNull"`
		Timestamp int     `json:"timestamp"`
		Value     string  `json:"value"`
	}

	var alt altStruct
	if err := json.Unmarshal(v, &alt); err == nil {
		a.Amount = alt.Amount
		a.Currency = alt.Currency
		a.IsNull = alt.IsNull
		a.Timestamp = alt.Timestamp

		val, err := strconv.ParseFloat(alt.Value, 64)
		if err != nil {
			a.Value = -1
		} else {
			a.Value = val
		}

		return nil
	}

	var inner AccountSummaryInner
	if err := json.Unmarshal(v, &inner); err != nil {
		return errors.Wrap(err, "failed to find appropriate struct to unmarshal account summary into")
	}

	a.Amount = inner.Amount
	a.Currency = inner.Currency
	a.IsNull = inner.IsNull
	a.Timestamp = inner.Timestamp
	a.Value = inner.Value

	return nil
}

/*
PortfolioAccounts - Gets portifolio accounts
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1accounts/get
*/
func (c *client) PortfolioAccounts() ([]PortfolioAccount, error) {
	resp, err := c.get(portfolioAccountsPath)
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

	portfolioAccounts := []PortfolioAccount{}
	if err := json.Unmarshal(v, &portfolioAccounts); err != nil {
		return nil, err
	}

	return portfolioAccounts, nil
}

/*
SubAccounts - Gets sub accounts
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1subaccounts/get
*/
func (c *client) SubAccounts() ([]SubAccount, error) {
	resp, err := c.get(subAccountsPath)
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

	subAccounts := []SubAccount{}
	if err := json.Unmarshal(v, &subAccounts); err != nil {
		return nil, err
	}

	return subAccounts, nil
}

/*
SubAccountsLarge - Gets sub accounts by page
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1subaccounts2/get
*/
func (c *client) SubAccountsLarge(page int) (*SubAccountsLarge, error) {
	resp, err := c.get(subAccountsLargePath, query{
		key:   "page",
		value: strconv.Itoa(page),
	})
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

	subAccountsLarge := &SubAccountsLarge{}
	if err := json.Unmarshal(v, subAccountsLarge); err != nil {
		return nil, err
	}

	return subAccountsLarge, nil
}

/*
AccountInformation - Gets information associated with an account
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1%7BaccountId%7D~1meta/get
*/
func (c *client) AccountInformation(accountID string) (*AccountInformation, error) {
	resp, err := c.get(
		substituteParam(accountInformationPath, param{
			key:   "accountId",
			value: accountID,
		}),
	)
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

	var accountInformation AccountInformation
	if err := json.Unmarshal(v, &accountInformation); err != nil {
		return nil, err
	}

	return &accountInformation, nil
}

/*
AccountSummary - Gets an accounts financial summary
Link: https://www.interactivebrokers.com/api/doc.html#tag/Portfolio/paths/~1portfolio~1%7BaccountId%7D~1summary/get
*/
func (c *client) AccountSummary(accountID string) (*AccountSummary, error) {
	resp, err := c.get(
		substituteParam(accountSummaryPath, param{
			key:   "accountId",
			value: accountID,
		}),
	)
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

	fmt.Println(string(v))

	var accountSummary AccountSummary
	if err := json.Unmarshal(v, &accountSummary); err != nil {
		return nil, err
	}

	return &accountSummary, nil
}
