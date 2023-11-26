package ibweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	newRequestFn = http.NewRequest
	readAllFn    = io.ReadAll
)

// Client - Client Portal Web API Interface
type Client interface {
	SetClient(httpClient *http.Client)

	// Contracts
	SearchContracts(input SearchContractsInput) ([]Contract, error)
	SearchStrikes(input SearchStrikesInput) (*SearchStrikes, error)
	SecurityDefinitionInfo(input SecurityDefinitionInfoInput) ([]SecurityDefinitionInfo, error)

	// Portfolio
	PortfolioAccounts() ([]PortfolioAccount, error)
	SubAccounts() ([]SubAccount, error)
	SubAccountsLarge(page int) (*SubAccountsLarge, error)
	AccountInformation(accountID string) (*AccountInformation, error)

	// Order
	PlaceOrders(accountID string, input PlaceOrdersInput) ([]PlaceOrders, error)
	PlaceOrderReply(replyID string, input PlaceOrderReplyInput) ([]PlaceOrders, error)
	CancelOrder(accountID, orderID string) (*CancelOrder, error)
	LiveOrders() (*LiveOrders, error)
	OrderStatus(orderID string) (*OrderStatus, error)

	//CCP
	PositionsByContractID(conID string) (*PositionsByConID, error)
}

type client struct {
	httpClient *http.Client
	url        string
	doFn       func(req *http.Request) (*http.Response, error)
}

// New - returns a new Client with the URL past
func New(url string) Client {
	c := http.DefaultClient
	return &client{
		httpClient: c,
		url:        url,
		doFn:       c.Do,
	}
}

// NewWithClient - retuns a new Client with the URL and *http.Client past.
func NewWithClient(httpClient *http.Client, url string) Client {
	return &client{
		httpClient: httpClient,
		url:        url,
		doFn:       httpClient.Do,
	}
}

// SetClient - sets the *http.Client
func (c *client) SetClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func substituteParam(path string, params ...param) string {
	for _, p := range params {
		key := "{" + p.key + "}"
		path = strings.ReplaceAll(path, key, p.value)
	}

	return path
}

type param struct {
	key   string
	value string
}

type query struct {
	key   string
	value string
}

func (c *client) get(path string, queries ...query) (*http.Response, error) {
	req, err := newRequestFn(
		http.MethodGet,
		fmt.Sprintf("%s/%s", c.url, path),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if queries != nil {
		q := req.URL.Query()
		for _, query := range queries {
			q.Add(query.key, query.value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return c.doFn(req)
}

func (c *client) post(path string, data interface{}) (*http.Response, error) {
	v, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := newRequestFn(
		http.MethodPost,
		fmt.Sprintf("%s/%s", c.url, path),
		bytes.NewBuffer(v),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.doFn(req)
}

func (c *client) delete(path string) (*http.Response, error) {
	req, err := newRequestFn(
		http.MethodDelete,
		fmt.Sprintf("%s/%s", c.url, path),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.doFn(req)
}
