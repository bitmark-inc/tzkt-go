package tzkt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TxsFormat struct {
	To      string `json:"to_"`
	Amount  string `json:"amount"`
	TokenID string `json:"token_id"`
}

type ParametersValue struct {
	From string      `json:"from_"`
	Txs  []TxsFormat `json:"txs"`
}

type TransactionParameter struct {
	EntryPoint string            `json:"entrypoint"`
	Value      []ParametersValue `json:"value"`
}

type Transaction struct {
	Block     string    `json:"block"`
	Target    Account   `json:"target"`
	Timestamp time.Time `json:"timestamp"`
	ID        uint64    `json:"id"`
	Hash      string    `json:"hash"`
}

type DetailedTransaction struct {
	Block     string               `json:"block"`
	Parameter TransactionParameter `json:"parameter"`
	Target    Account              `json:"target"`
	Timestamp time.Time            `json:"timestamp"`
	ID        uint64               `json:"id"`
	Hash      string               `json:"hash"`
}

// GetTransactionByTx gets transaction details from a specific Tx
// confirmed => true; failed => false; pending => nil
func (c *TZKT) GetTransactionStatusByTx(hash string) (*bool, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.endpoint,
		Path:   fmt.Sprintf("%s/%s/%s", "/v1/operations/transactions", hash, "status"),
	}

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to get transaction status: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	status, err := strconv.ParseBool(string(b))
	if err != nil {
		return nil, nil
	}
	return &status, nil
}

// GetTransactionByTx gets transaction details from a specific Tx
func (c *TZKT) GetTransactionByTx(hash string) ([]DetailedTransaction, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.endpoint,
		Path:   fmt.Sprintf("%s/%s", "/v1/operations/transactions", hash),
	}

	var transactionDetails []DetailedTransaction

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

func (c *TZKT) GetTransaction(id uint64) (Transaction, error) {
	v := url.Values{
		"id": []string{fmt.Sprintf("%d", id)},
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/operations/transactions",
		RawQuery: v.Encode(),
	}

	var txs []Transaction

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Transaction{}, err
	}
	if err := c.request(req, &txs); err != nil {
		return Transaction{}, err
	}

	if len(txs) == 0 {
		return Transaction{}, fmt.Errorf("transaction not found")
	}
	return txs[0], nil
}
