package ibweb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	placeOrdersPath     = "v1/api/iserver/account/{accountId}/orders"
	cancelOrderPath     = "v1/api/iserver/account/{accountId}/order/{orderId}"
	placeOrderReplyPath = "v1/api/iserver/reply/{replyid}"
	liveOrdersPath      = "v1/api/iserver/account/orders"
	orderStatusPath     = "v1/api/iserver/account/order/status/{orderId}"
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

/*
PlaceOrders -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1%7BaccountId%7D~1orders/post
*/
type PlaceOrders struct {
	OrderID          string   `json:"order_id"`
	OrderStatus      string   `json:"order_status"`
	EncryptedMessage string   `json:"encrypt_message"`
	ID               string   `json:"id"`
	Message          []string `json:"message"`
}

/*
PlaceOrderReplyInput - Reply to an interactive brokers order confirmation check
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1reply~1%7Breplyid%7D/post
*/
type PlaceOrderReplyInput struct {
	Confirmed bool `json:"confirmed"`
}

/*
PlaceOrderReplyInput -
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1reply~1%7Breplyid%7D/post
*/
type PlaceOrderReply struct {
	OrderID      string `json:"order_id"`
	OrderStatus  string `json:"order_status"`
	LocalOrderID string `json:"local_order_id"`
}

/*
CancelOrder - Cancel an order placed
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1%7BaccountId%7D~1order~1%7BorderId%7D/delete
*/
type CancelOrder struct {
	OrderID int    `json:"order_id"`
	Msg     string `json:"msg"`
	Conid   int    `json:"conid"`
	Account string `json:"account,omitempty"`
}

type LiveOrders struct {
	Filters []string `json:"filters"`
	Orders  []struct {
		Acct               string  `json:"acct"`
		Conidex            string  `json:"conidex"`
		Conid              int     `json:"conid"`
		OrderID            int     `json:"orderId"`
		CashCcy            string  `json:"cashCcy"`
		SizeAndFills       string  `json:"sizeAndFills"`
		OrderDesc          string  `json:"orderDesc"`
		Description1       string  `json:"description1"`
		Ticker             string  `json:"ticker"`
		SecType            string  `json:"secType"`
		ListingExchange    string  `json:"listingExchange"`
		RemainingQuantity  float64 `json:"remainingQuantity"`
		FilledQuantity     float64 `json:"filledQuantity"`
		CompanyName        string  `json:"companyName"`
		Status             string  `json:"status"`
		OrigOrderType      string  `json:"origOrderType"`
		SupportsTaxOpt     string  `json:"supportsTaxOpt"`
		LastExecutionTime  string  `json:"lastExecutionTime"`
		LastExecutionTimeR int     `json:"lastExecutionTime_r"`
		OrderType          string  `json:"orderType"`
		OrderRef           string  `json:"order_ref"`
		Side               string  `json:"side"`
		TimeInForce        string  `json:"timeInForce"`
		Price              int     `json:"price"`
		BgColor            string  `json:"bgColor"`
		FgColor            string  `json:"fgColor"`
	} `json:"orders"`
	Snapshot bool `json:"snapshot"`
}

/*
OrderStatus - Gets the status of an order
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1order~1status~1%7BorderId%7D/get
*/
type OrderStatus struct {
	SubType                      string `json:"sub_type"`
	RequestID                    string `json:"request_id"`
	OrderID                      int    `json:"order_id"`
	Conidex                      string `json:"conidex"`
	Symbol                       string `json:"symbol"`
	Side                         string `json:"side"`
	ContractDescription1         string `json:"contract_description_1"`
	ListingExchange              string `json:"listing_exchange"`
	OptionAcct                   string `json:"option_acct"`
	CompanyName                  string `json:"company_name"`
	Size                         string `json:"size"`
	TotalSize                    string `json:"total_size"`
	Currency                     string `json:"currency"`
	Account                      string `json:"account"`
	OrderType                    string `json:"order_type"`
	LimitPrice                   string `json:"limit_price"`
	StopPrice                    string `json:"stop_price"`
	CumFill                      string `json:"cum_fill"`
	OrderStatus                  string `json:"order_status"`
	OrderStatusDescription       string `json:"order_status_description"`
	Tif                          string `json:"tif"`
	FgColor                      string `json:"fg_color"`
	BgColor                      string `json:"bg_color"`
	OrderNotEditable             bool   `json:"order_not_editable"`
	EditableFields               string `json:"editable_fields"`
	CannotCancelOrder            bool   `json:"cannot_cancel_order"`
	OutsideRth                   bool   `json:"outside_rth"`
	DeactivateOrder              bool   `json:"deactivate_order"`
	UsePriceMgmtAlgo             bool   `json:"use_price_mgmt_algo"`
	SecType                      string `json:"sec_type"`
	AvailableChartPeriods        string `json:"available_chart_periods"`
	OrderDescription             string `json:"order_description"`
	OrderDescriptionWithContract string `json:"order_description_with_contract"`
	AlertActive                  int    `json:"alert_active"`
	ChildOrderType               string `json:"child_order_type"`
	SizeAndFills                 string `json:"size_and_fills"`
	ExitStrategyDisplayPrice     string `json:"exit_strategy_display_price"`
	ExitStrategyChartDescription string `json:"exit_strategy_chart_description"`
	ExitStrategyToolAvailability int    `json:"exit_strategy_tool_availability"`
	AllowedDuplicateOpposite     bool   `json:"allowed_duplicate_opposite"`
	OrderTime                    string `json:"order_time"`
	OcaGroupID                   string `json:"oca_group_id"`
}

/*
PlaceOrders - Places orders
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1%7BaccountId%7D~1orders/post
*/
func (c *client) PlaceOrders(accountID string, input PlaceOrdersInput) ([]PlaceOrders, error) {
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

	var orderResp []PlaceOrders
	if err := json.Unmarshal(v, &orderResp); err != nil {
		return nil, err
	}

	return orderResp, nil
}

/*
PlaceOrderReply - Reply to an interactive brokers order confirmation check
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1reply~1%7Breplyid%7D/post
*/
func (c *client) PlaceOrderReply(replyID string, input PlaceOrderReplyInput) ([]PlaceOrders, error) {
	resp, err := c.post(substituteParam(placeOrderReplyPath, param{
		key:   "replyid",
		value: replyID,
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

	var orderResp []PlaceOrders
	if err := json.Unmarshal(v, &orderResp); err != nil {
		return nil, err
	}

	return orderResp, nil
}

/*
CancelOrder - Cancel an order placed
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1%7BaccountId%7D~1order~1%7BorderId%7D/delete
*/
func (c *client) CancelOrder(accountID, orderID string) (*CancelOrder, error) {
	resp, err := c.delete(substituteParam(cancelOrderPath,
		param{
			key:   "accountId",
			value: accountID,
		},
		param{
			key:   "orderId",
			value: orderID,
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

	var cancelOrder CancelOrder
	if err := json.Unmarshal(v, &cancelOrder); err != nil {
		return nil, err
	}

	return &cancelOrder, nil
}

/*
LiveOrders - Gets all live orders
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order
*/
func (c *client) LiveOrders() (*LiveOrders, error) {
	resp, err := c.get(liveOrdersPath)
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

	var liveOrders LiveOrders
	if err := json.Unmarshal(v, &liveOrders); err != nil {
		return nil, err
	}

	return &liveOrders, nil
}

/*
OrderStatus - Gets the status of an order
Link: https://www.interactivebrokers.com/api/doc.html#tag/Order/paths/~1iserver~1account~1order~1status~1%7BorderId%7D/get
*/
func (c *client) OrderStatus(orderID string) (*OrderStatus, error) {
	resp, err := c.get(substituteParam(orderStatusPath,
		param{
			key:   "orderId",
			value: orderID,
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
	fmt.Println(string(v))

	var orderStatus OrderStatus
	if err := json.Unmarshal(v, &orderStatus); err != nil {
		return nil, err
	}

	return &orderStatus, nil
}
