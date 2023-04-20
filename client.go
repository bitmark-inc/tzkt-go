package tzkt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrTooManyRequest = fmt.Errorf("too many requests")

type TZKT struct {
	endpoint string

	client *http.Client
}

func New(network string) *TZKT {
	endpoint := "api.mainnet.tzkt.io"
	if network == "testnet" {
		endpoint = "api.ghostnet.tzkt.io"
	}

	return &TZKT{
		client: &http.Client{
			Timeout: time.Minute,
		},
		endpoint: endpoint,
	}
}

func (c *TZKT) request(req *http.Request, responseData interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// close the body only when we return an error
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			return ErrTooManyRequest
		}

		errResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("tzkt api error: %s", errResp)
	}

	err = json.NewDecoder(resp.Body).Decode(&responseData)

	return err
}
