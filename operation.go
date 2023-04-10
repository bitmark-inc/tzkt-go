package tzkt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetTransactionByTx gets transaction details from a specific Tx
func (c *TZKT) GetTransactionByTx(hash string) ([]TransactionDetails, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.endpoint,
		Path:   fmt.Sprintf("%s/%s", "/v1/operations/transactions", hash),
	}

	var transactionDetails []TransactionDetails

	resp, err := c.client.Get(u.String())
	if err != nil {
		return transactionDetails, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&transactionDetails); err != nil {
		return transactionDetails, err
	}

	return transactionDetails, nil
}

func (c *TZKT) GetTransaction(id uint64) (TransactionDetails, error) {
	v := url.Values{
		"id": []string{fmt.Sprintf("%d", id)},
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/operations/transactions",
		RawQuery: v.Encode(),
	}

	var txs []TransactionDetails

	req, _ := http.NewRequest("GET", u.String(), nil)
	if err := c.request(req, &txs); err != nil {
		return TransactionDetails{}, err
	}

	if len(txs) == 0 {
		return TransactionDetails{}, fmt.Errorf("transaction not found")
	}
	return txs[0], nil
}
