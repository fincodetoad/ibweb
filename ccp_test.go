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

	portfolioAccounts, err := c.PortfolioAccounts()
	assert.Nil(t, err)
	assert.Greater(t, len(portfolioAccounts), 0)

	if !t.Failed() {
		_, err := c.PositionByContractID(portfolioAccounts[0].AccountID, "659248794")
		assert.Nil(t, err)
	}
}
