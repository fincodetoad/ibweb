package ibweb

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestPlaceOrderReplCancelAndLiveIntegration(t *testing.T) {
	c := NewWithClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}, "https://127.0.0.1:5555")

	portfolioAccounts, err := c.PortfolioAccounts()
	assert.Nil(t, err)
	assert.Greater(t, len(portfolioAccounts), 0)

	var orderID string
	if !t.Failed() {
		orders, err := c.PlaceOrders(portfolioAccounts[0].AccountID, PlaceOrdersInput{
			Orders: []Order{
				{
					AcctID:    portfolioAccounts[0].AccountID,
					Conid:     659248794,
					Quantity:  1,
					OrderType: Market,
					Side:      Buy,
					Tif:       GoodTillCanceled,
				},
			},
		})
		assert.Nil(t, err)
		if !assert.Greater(t, len(orders), 0) {
			t.FailNow()
		}

		if orders[0].ID != "" {
			replies, err := c.PlaceOrderReply(orders[0].ID, PlaceOrderReplyInput{Confirmed: true})
			if !assert.Nil(t, err) {
				t.FailNow()
			}
			if !assert.Greater(t, len(replies), 0) {
				t.FailNow()
			}

			orderID = replies[0].OrderID
		} else {
			orderID = orders[0].OrderID
		}
	}

	if !t.Failed() {
		_, err := c.LiveOrders()
		assert.Nil(t, err)
	}

	if !t.Failed() {
		_, err := c.CancelOrder(portfolioAccounts[0].AccountID, orderID)
		assert.Nil(t, err)
	}
}

func TestPlaceOrdersUnit(t *testing.T) {
	type input struct {
		handler   func(req *http.Request) (*http.Response, error)
		readAllFn func(r io.Reader) ([]byte, error)
	}

	type want struct {
		wantErr         bool
		wantErrContains string
	}

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to post",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("failed to post orders")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to post orders",
			},
		},
		{
			"handles unexpected status code",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(500, "failed"), nil
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid status code",
			},
		},
		{
			"handles failure to read response body",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, ""), nil
				},
				readAllFn: func(r io.Reader) ([]byte, error) {
					return nil, errors.New("failed to read order response")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read order response",
			},
		},
		{
			"handles failure to unmarshal response",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, "garbage"), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid character",
			},
		},
		{
			"is successful",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					v, err := os.ReadFile("./testdata/placeorders.json")
					if !assert.Nil(t, err) {
						t.FailNow()
					}

					return httpmock.NewBytesResponse(200, v), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr: false,
			},
		},
	}

	for _, tc := range tests {
		httpmock.Activate()
		readAllFn = tc.input.readAllFn

		httpmock.RegisterResponder(http.MethodPost, "http://127.0.0.1:5555/v1/api/iserver/account/DU777777/orders", tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.PlaceOrders("DU777777", PlaceOrdersInput{})
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestPlaceOrderReplyUnit(t *testing.T) {
	type input struct {
		handler   func(req *http.Request) (*http.Response, error)
		readAllFn func(r io.Reader) ([]byte, error)
	}

	type want struct {
		wantErr         bool
		wantErrContains string
	}

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to post",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("failed to post order reply")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to post order reply",
			},
		},
		{
			"handles unexpected status code",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(500, "failed"), nil
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid status code",
			},
		},
		{
			"handles failure to read response body",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, ""), nil
				},
				readAllFn: func(r io.Reader) ([]byte, error) {
					return nil, errors.New("failed to read order reply response")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read order reply response",
			},
		},
		{
			"handles failure to unmarshal response",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, "garbage"), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid character",
			},
		},
		{
			"is successful",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					v, err := os.ReadFile("./testdata/order_reply.json")
					if !assert.Nil(t, err) {
						t.FailNow()
					}

					return httpmock.NewBytesResponse(200, v), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr: false,
			},
		},
	}

	for _, tc := range tests {
		httpmock.Activate()
		readAllFn = tc.input.readAllFn

		httpmock.RegisterResponder(http.MethodPost, "http://127.0.0.1:5555/v1/api/iserver/reply/888888", tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.PlaceOrderReply("888888", PlaceOrderReplyInput{})
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestCancelOrderUnit(t *testing.T) {
	type input struct {
		handler   func(req *http.Request) (*http.Response, error)
		readAllFn func(r io.Reader) ([]byte, error)
	}

	type want struct {
		wantErr         bool
		wantErrContains string
	}

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to post",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("failed to cancel orders")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to cancel orders",
			},
		},
		{
			"handles unexpected status code",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(500, "failed"), nil
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid status code",
			},
		},
		{
			"handles failure to read response body",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, ""), nil
				},
				readAllFn: func(r io.Reader) ([]byte, error) {
					return nil, errors.New("failed to read order response")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read order response",
			},
		},
		{
			"handles failure to unmarshal response",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, "garbage"), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid character",
			},
		},
		{
			"is successful",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					v, err := os.ReadFile("./testdata/cancel_order.json")
					if !assert.Nil(t, err) {
						t.FailNow()
					}

					return httpmock.NewBytesResponse(200, v), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr: false,
			},
		},
	}

	for _, tc := range tests {
		httpmock.Activate()
		readAllFn = tc.input.readAllFn

		httpmock.RegisterResponder(http.MethodDelete, "http://127.0.0.1:5555/v1/api/iserver/account/DU777777/order/22345544", tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.CancelOrder("DU777777", "22345544")
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestLiveOrdersUnit(t *testing.T) {
	type input struct {
		handler   func(req *http.Request) (*http.Response, error)
		readAllFn func(r io.Reader) ([]byte, error)
	}

	type want struct {
		wantErr         bool
		wantErrContains string
	}

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to get",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("failed to get live orders")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get live orders",
			},
		},
		{
			"handles unexpected status code",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(500, "failed"), nil
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid status code",
			},
		},
		{
			"handles failure to read response body",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, ""), nil
				},
				readAllFn: func(r io.Reader) ([]byte, error) {
					return nil, errors.New("failed to get live orders")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get live orders",
			},
		},
		{
			"handles failure to unmarshal response",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(200, "garbage"), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr:         true,
				wantErrContains: "invalid character",
			},
		},
		{
			"is successful",
			input{
				handler: func(req *http.Request) (*http.Response, error) {
					v, err := os.ReadFile("./testdata/live_orders.json")
					if !assert.Nil(t, err) {
						t.FailNow()
					}

					return httpmock.NewBytesResponse(200, v), nil
				},
				readAllFn: io.ReadAll,
			},
			want{
				wantErr: false,
			},
		},
	}

	for _, tc := range tests {
		httpmock.Activate()
		readAllFn = tc.input.readAllFn

		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("http://127.0.0.1:5555/%s", liveOrdersPath), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.LiveOrders()
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}
