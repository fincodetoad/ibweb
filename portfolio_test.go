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

func TestPortfolioAccountsIntegration(t *testing.T) {
	c := NewWithClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}, "https://127.0.0.1:5555")

	_, err := c.PortfolioAccounts()
	assert.Nil(t, err)
}

func TestPortfolioAccountsUnit(t *testing.T) {
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
					return nil, errors.New("failed to get portfolio accounts")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get portfolio accounts",
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
					return nil, errors.New("failed to read portfolio accounts body")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read portfolio accounts body",
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
					v, err := os.ReadFile("./testdata/portfolio_accounts.json")
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

		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("http://127.0.0.1:5555/%s", portfolioAccountsPath), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.PortfolioAccounts()
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestSubAccountsIntegration(t *testing.T) {
	c := NewWithClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}, "https://127.0.0.1:5555")

	_, err := c.SubAccounts()
	assert.Nil(t, err)
}

func TestSubAccountsUnit(t *testing.T) {
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
					return nil, errors.New("failed to get subaccounts")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get subaccounts",
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
					return nil, errors.New("failed to read subaccounts body")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read subaccounts body",
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
					v, err := os.ReadFile("./testdata/subaccounts.json")
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

		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("http://127.0.0.1:5555/%s", subAccountsPath), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.SubAccounts()
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestSubAccountsLargeIntegration(t *testing.T) {
	c := NewWithClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}, "https://127.0.0.1:5555")

	_, err := c.SubAccountsLarge(0)
	assert.Nil(t, err)
}

func TestSubAccountsLargeUnit(t *testing.T) {
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
					return nil, errors.New("failed to get subaccounts large")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get subaccounts large",
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
					return nil, errors.New("failed to read subaccounts large body")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read subaccounts large body",
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
					v, err := os.ReadFile("./testdata/subaccounts_large.json")
					if !assert.Nil(t, err) {
						t.FailNow()
					}

					_, ok := req.URL.Query()["page"]
					assert.True(t, ok)

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

		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("http://127.0.0.1:5555/%s", subAccountsLargePath), tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.SubAccountsLarge(0)
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}

func TestAccountInformationIntegration(t *testing.T) {
	c := NewWithClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}, "https://127.0.0.1:5555")

	portfolioAccounts, err := c.PortfolioAccounts()
	assert.Nil(t, err)
	assert.Greater(t, len(portfolioAccounts), 0)

	_, err = c.AccountInformation(portfolioAccounts[0].AccountID)
	assert.Nil(t, err)
	t.Fail()
}

func TestAccountInformationUnit(t *testing.T) {
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
					return nil, errors.New("failed to get account information")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to get account information",
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
					return nil, errors.New("failed to read account information")
				},
			},
			want{
				wantErr:         true,
				wantErrContains: "failed to read account information",
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
					v, err := os.ReadFile("./testdata/account_information.json")
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

		httpmock.RegisterResponder(http.MethodGet, "http://127.0.0.1:5555/v1/api/portfolio/DU7777777/meta", tc.input.handler)

		c := New("http://127.0.0.1:5555")
		_, err := c.AccountInformation("DU7777777")
		assertError(t, tc.want.wantErr, tc.want.wantErrContains, err)

		httpmock.DeactivateAndReset()
	}
}
