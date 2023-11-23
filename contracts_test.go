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

func TestSearchContractsIntegration(t *testing.T) {
	c := NewWithClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}, "https://127.0.0.1:5555")

	contracts, err := c.SearchContracts(SearchContractsInput{
		Symbol:  "AAPL",
		SecType: Options,
	})
	assert.Nil(t, err)
	assert.Greater(t, len(contracts), 0)
}

func TestSearchContractsUnit(t *testing.T) {
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
					return nil, errors.New("failed to post contracts search")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to post contracts search",
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
					return nil, errors.New("failed to read contracts")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read contracts",
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
					v, err := os.ReadFile("./testdata/contracts.json")
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

		httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("http://127.0.0.1:5555/%s", searchContractsPath), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.SearchContracts(SearchContractsInput{})
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestSearchStrikesIntegration(t *testing.T) {
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
	}
}

func TestSearchStrikesUnit(t *testing.T) {
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
					return nil, errors.New("failed to get strikes")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get strikes",
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
					return nil, errors.New("failed to read strikes")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read strikes",
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
					v, err := os.ReadFile("./testdata/search_strikes.json")
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

		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("http://127.0.0.1:5555/%s", searchStrikesPath), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.SearchStrikes(SearchStrikesInput{})
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestSecurityDefinitionInfoIntegration(t *testing.T) {
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

		if !t.Failed() {
			// now := time.Now()
			secDefIndo, err := c.SecurityDefinitionInfo(SecurityDefinitionInfoInput{
				ConID:   contracts[0].Conid,
				SecType: Options,
				Month:   "DEC23",
				Strike:  strikes.Call[0],
			})
			assert.Nil(t, err)
			assert.Greater(t, len(secDefIndo), 0)
		}
	}
}

func TestSecurityDefinitionInfoUnit(t *testing.T) {
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
					return nil, errors.New("failed to get sec def info")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get sec def info",
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
					return nil, errors.New("failed to read sec def info")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read sec def info",
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
					v, err := os.ReadFile("./testdata/secdefinfo.json")
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

		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("http://127.0.0.1:5555/%s", secDefInfoPath), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.SecurityDefinitionInfo(SecurityDefinitionInfoInput{})
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}
