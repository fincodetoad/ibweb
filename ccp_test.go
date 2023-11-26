package ibweb

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionsByContractIDIntegration(t *testing.T) {
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
		_, err := c.PositionsByContractID(contracts[0].Conid)
		assert.Nil(t, err)
	}
}
