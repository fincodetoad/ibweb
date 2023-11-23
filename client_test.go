package ibweb

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	c := http.DefaultClient

	type input struct {
		newRequestFn func(method string, url string, body io.Reader) (*http.Request, error)
		doFn         func(req *http.Request) (*http.Response, error)
		queries      []query
	}

	type want struct {
		wantErr        bool
		wantErrMessage string
		assertions     []func(t *testing.T, r *http.Request)
	}

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to create new request",
			input{
				newRequestFn: func(method, url string, body io.Reader) (*http.Request, error) {
					return nil, errors.New("failed to create new http request")
				},
			},
			want{
				wantErr:        true,
				wantErrMessage: "failed to create new http request",
			},
		},
		{
			"handles http do failure",
			input{
				newRequestFn: http.NewRequest,
				doFn: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("failed to do")
				},
			},
			want{
				wantErr:        true,
				wantErrMessage: "failed to do",
			},
		},
		{
			"is successful",
			input{
				newRequestFn: http.NewRequest,
				doFn:         c.Do,
				queries: []query{
					{
						key:   "key1",
						value: "val1",
					},
					{
						key:   "key2",
						value: "val2",
					},
				},
			},
			want{
				wantErr: false,
				assertions: []func(t *testing.T, r *http.Request){
					func(t *testing.T, r *http.Request) {
						vals := r.URL.Query()
						val, ok := vals["key1"]
						if !assert.True(t, ok) {
							return
						}
						assert.Equal(t, []string{"val1"}, val)

						val, ok = vals["key2"]
						if !assert.True(t, ok) {
							return
						}
						assert.Equal(t, []string{"val2"}, val)

					},
				},
			},
		},
	}

	for _, tc := range tests {
		httpmock.Activate()
		newRequestFn = tc.input.newRequestFn

		c := client{
			httpClient: c,
			url:        "http://127.0.0.1:5555",
			doFn:       tc.input.doFn,
		}

		httpmock.RegisterResponder("GET", "http://127.0.0.1:5555/test",
			func(req *http.Request) (*http.Response, error) {
				for _, assertion := range tc.want.assertions {
					assertion(t, req)
				}

				return httpmock.NewStringResponse(200, ""), nil
			})

		_, err := c.get("test", tc.input.queries...)
		assertError(t, tc.want.wantErr, tc.want.wantErrMessage, err)

		httpmock.DeactivateAndReset()
	}
}

func assertError(t *testing.T, wantError bool, message string, err error) bool {
	if wantError {
		if !assert.NotNil(t, err) {
			return false
		}

		if !assert.Contains(t, err.Error(), message) {
			return false
		}
	} else {
		if !assert.Nil(t, err) {
			return false
		}
	}

	return true
}
