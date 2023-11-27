package ibweb

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMarketDataHistoryIntegration(t *testing.T) {
	c := NewWithClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}, "https://127.0.0.1:5555")

	contracts, err := c.SearchContracts(SearchContractsInput{
		Symbol:  "AAPL",
		SecType: Options,
	})
	assert.Nil(t, err)
	assert.Greater(t, len(contracts), 0)
	if !t.Failed() {
		// now := time.Now()
		strikes, err := c.SearchStrikes(SearchStrikesInput{
			ConID:   contracts[0].Conid,
			SecType: Options,
			Month:   "DEC23",
		})
		assert.Nil(t, err)
		assert.Greater(t, len(strikes.Call), 0)
		assert.Greater(t, len(strikes.Put), 0)

		conID := ""
		if !t.Failed() {
			// now := time.Now()
			secDefInfo, err := c.SecurityDefinitionInfo(SecurityDefinitionInfoInput{
				ConID:   contracts[0].Conid,
				SecType: Options,
				Month:   "DEC23",
				Strike:  strikes.Call[0],
			})
			assert.Nil(t, err)
			assert.Greater(t, len(secDefInfo), 0)
			conID = strconv.Itoa(secDefInfo[0].Conid)

		}

		if !t.Failed() {
			_, err := c.MarketDataHistory(MarketDataHistoryInput{ConID: conID})
			assert.Nil(t, err)
		}
	}
}

func TestMarketDataHistoryUnit(t *testing.T) {
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
					return nil, errors.New("failed to get market data history")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get market data history",
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
					return nil, errors.New("failed to read market data history")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read market data history",
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
					v, err := os.ReadFile("./testdata/market_data_history.json")
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

		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("http://127.0.0.1:5555/%s", marketDataHistory), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.MarketDataHistory(MarketDataHistoryInput{})
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}
