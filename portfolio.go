package ibweb

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	portfolioAccountsPath  = "v1/api/portfolio/accounts"
	subAccountsPath        = "v1/api/portfolio/subaccounts"
	subAccountsLargePath   = "v1/api/portfolio/subaccounts2"
	accountInformationPath = "v1/api/portfolio/{accountId}/meta"
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
