package tzkt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

var ErrTooManyRequest = fmt.Errorf("too many requests")

type TZKT struct {
	endpoint string

	client *http.Client
}

type NullableInt int64

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

type FormatDimensions struct {
	Unit  string `json:"unit"`
	Value string `json:"value"`
}

type FileFormat struct {
	URI        string           `json:"uri"`
	FileName   string           `json:"fileName,omitempty"`
	FileSize   int              `json:"fileSize,string"`
	MIMEType   MIMEFormat       `json:"mimeType"`
	Dimensions FormatDimensions `json:"dimensions,omitempty"`
}

type MIMEFormat string

func (m *MIMEFormat) UnmarshalJSON(data []byte) error {
	if data[0] == 91 {
		data = bytes.Trim(data, "[]")
	}

	return json.Unmarshal(data, (*string)(m))
}

type FileFormats []FileFormat

func (t *NullableInt) UnmarshalJSON(data []byte) error {
	var num int64

	err := json.Unmarshal(data, &num)
	if err != nil {
		*t = NullableInt(-1)

		return nil
	}

	*t = NullableInt(num)

	return nil
}

func (f *FileFormats) UnmarshalJSON(data []byte) error {
	type formats FileFormats

	switch data[0] {
	case 34:
		d1 := bytes.ReplaceAll(bytes.Trim(data, `"`), []byte{92, 117, 48, 48, 50, 50}, []byte{34})
		d := bytes.ReplaceAll(d1, []byte{92, 34}, []byte{34})

		if err := json.Unmarshal(d, (*formats)(f)); err != nil {
			return err
		}
	case 123: // If the "formats" is not an array
		d := append([]byte{91}, data...)
		if data[len(data)-1] != 93 {
			d = append(d, []byte{93}...)
		}
		if err := json.Unmarshal(d, (*formats)(f)); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(data, (*formats)(f)); err != nil {
			return err
		}
	}

	return nil
}

type FileCreators []string

func (c *FileCreators) UnmarshalJSON(data []byte) error {
	type creators FileCreators

	switch data[0] {
	case 34:
		d1 := bytes.ReplaceAll(bytes.Trim(data, `"`), []byte{92, 117, 48, 48, 50, 50}, []byte{34})
		d := bytes.ReplaceAll(d1, []byte{92, 34}, []byte{34})

		if err := json.Unmarshal(d, (*creators)(c)); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(data, (*creators)(c)); err != nil {
			return err
		}
	}

	return nil
}

type TokenID struct {
	big.Int
}

func (b TokenID) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *TokenID) UnmarshalJSON(p []byte) error {
	s := string(p)

	if s == "null" {
		return fmt.Errorf("invalid token id: %s", p)
	}

	z, ok := big.NewInt(0).SetString(strings.Trim(s, `"`), 0)
	if !ok {
		return fmt.Errorf("invalid token id: %s", p)
	}

	b.Int = *z
	return nil
}

type Account struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

type Token struct {
	Contract    Account        `json:"contract"`
	ID          TokenID        `json:"tokenId"`
	Standard    string         `json:"standard"`
	TotalSupply NullableInt    `json:"totalSupply,string"`
	Timestamp   time.Time      `json:"firstTime"`
	Metadata    *TokenMetadata `json:"metadata,omitempty"`
}

type OwnedToken struct {
	Token    Token       `json:"token"`
	Balance  NullableInt `json:"balance,string"`
	LastTime time.Time   `json:"lastTime"`
}

type TokenMetadata struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Symbol       string       `json:"symbol"`
	RightURI     string       `json:"rightUri"`
	ArtifactURI  string       `json:"artifactUri"`
	DisplayURI   string       `json:"displayUri"`
	ThumbnailURI string       `json:"thumbnailUri"`
	Publishers   []string     `json:"publishers"`
	Creators     FileCreators `json:"creators"`
	Formats      FileFormats  `json:"formats"`

	ArtworkMetadata map[string]interface{} `json:"artworkMetadata"`
}

type TokenTransfer struct {
	Timestamp     time.Time `json:"timestamp"`
	Level         uint64    `json:"level"`
	TransactionID uint64    `json:"transactionId"`
	From          *Account  `json:"from"`
	To            Account   `json:"to"`
}

type TokenOwner struct {
	Address  string    `json:"address"`
	Balance  int64     `json:"balance,string"`
	LastTime time.Time `json:"lastTime"`
}

type TransactionDetails struct {
	Block     string               `json:"block"`
	Parameter TransactionParameter `json:"parameter"`
	Target    Account              `json:"target"`
	Timestamp time.Time            `json:"timestamp"`
	ID        uint64               `json:"id"`
	Hash      string               `json:"hash"`
}

type TransactionParameter struct {
	EntryPoint string            `json:"entrypoint"`
	Value      []ParametersValue `json:"value"`
}

type ParametersValue struct {
	From string      `json:"from_"`
	Txs  []TxsFormat `json:"txs"`
}

type TxsFormat struct {
	To      string `json:"to_"`
	Amount  string `json:"amount"`
	TokenID string `json:"token_id"`
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
