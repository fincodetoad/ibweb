package ibweb

import (
	"encoding/json"
	"net/http"
)

const (
	placeOrdersPath = "v1/api/iserver/account/{accountId}/orders"
)

// TimeInForce - time for order to execute
type TimeInForce string

const (
	GoodTillCanceled   = "GTC"
	OpenPriceGuarantee = "OPG"
	Dat                = "DAY"
	ImediateOrCancel   = "IOC"
)

// OrderSide - buy or sell
type OrderSide string

const (
	Buy  OrderSide = "BUY"
	Sell OrderSide = "SELL"
)

// OrderType - supported order types
type OrderType string

const (
	Limit      OrderType = "LMT"
	Market     OrderType = "MKT"
	Stop       OrderType = "STP"
	StopLimit  OrderType = "STOP_LIMIT"
	MidPrice   OrderType = "MIDPRICE"
	Trail      OrderType = "TRAIL"
	TrailLimit OrderType = "TRAILLMT"
)

/*
PlaceOrdersInput -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1%7BaccountId%7D~1orders/post
*/
type PlaceOrdersInput struct {
	Orders []Order `json:"orders"`
}

/*
Order -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1%7BaccountId%7D~1orders/post
*/
type Order struct {
	AcctID             string      `json:"acctId,omitempty"`
	Conid              int         `json:"conid,omitempty"`
	Conidex            string      `json:"conidex,omitempty"`
	SecType            string      `json:"secType,omitempty"`
	COID               string      `json:"cOID,omitempty"`
	ParentID           string      `json:"parentId,omitempty"`
	OrderType          OrderType   `json:"orderType,omitempty"`
	ListingExchange    string      `json:"listingExchange,omitempty"`
	IsSingleGroup      bool        `json:"isSingleGroup,omitempty"`
	OutsideRTH         bool        `json:"outsideRTH,omitempty"`
	Price              int         `json:"price,omitempty"`
	AuxPrice           interface{} `json:"auxPrice,omitempty"`
	Side               OrderSide   `json:"side,omitempty"`
	Ticker             string      `json:"ticker,omitempty"`
	Tif                TimeInForce `json:"tif,omitempty"`
	TrailingAmt        int         `json:"trailingAmt,omitempty"`
	TrailingType       string      `json:"trailingType,omitempty"`
	Referrer           string      `json:"referrer,omitempty"`
	Quantity           int         `json:"quantity,omitempty"`
	CashQty            int         `json:"cashQty,omitempty"`
	FxQty              int         `json:"fxQty,omitempty"`
	UseAdaptive        bool        `json:"useAdaptive,omitempty"`
	IsCcyConv          bool        `json:"isCcyConv,omitempty"`
	AllocationMethod   string      `json:"allocationMethod,omitempty"`
	Strategy           string      `json:"strategy,omitempty"`
	StrategyParameters struct {
	} `json:"strategyParameters,omitempty"`
}

type OrderPlacedResp struct {
	OrderID          string   `json:"order_id"`
	OrderStatus      string   `json:"order_status"`
	EncryptedMessage string   `json:"encrypt_message"`
	ID               string   `json:"id"`
	Message          []string `json:"message"`
}

/*
PlaceOrders - Places orders
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1%7BaccountId%7D~1orders/post
*/
func (c *client) PlaceOrders(accountID string, input PlaceOrdersInput) ([]OrderPlacedResp, error) {
	resp, err := c.post(substituteParam(placeOrdersPath, param{
		key:   "accountId",
		value: accountID,
	}), input)
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

	var orderResp []OrderPlacedResp
	if err := json.Unmarshal(v, &orderResp); err != nil {
		return nil, err
	}

	return orderResp, nil
}
