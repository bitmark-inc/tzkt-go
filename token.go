package tzkt

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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

type Token struct {
	Contract    Account        `json:"contract"`
	ID          TokenID        `json:"tokenId"`
	Standard    string         `json:"standard"`
	TotalSupply NullableInt    `json:"totalSupply,string"`
	Timestamp   time.Time      `json:"firstTime"`
	LastTime    time.Time      `json:"lastTime"`
	Metadata    *TokenMetadata `json:"metadata,omitempty"`
}

type OwnedToken struct {
	Token     Token       `json:"token"`
	Balance   NullableInt `json:"balance,string"`
	FirstTime time.Time   `json:"firstTime"`
	LastTime  time.Time   `json:"lastTime"`
}

type TokenMetadata struct {
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	Symbol          string       `json:"symbol"`
	RightURI        string       `json:"rightUri"`
	ArtifactURI     string       `json:"artifactUri"`
	DisplayURI      string       `json:"displayUri"`
	IsBooleanAmount StringBool   `json:"isBooleanAmount"`
	ThumbnailURI    string       `json:"thumbnailUri"`
	Publishers      []string     `json:"publishers"`
	Minter          string       `json:"minter"`
	Creators        FileCreators `json:"creators"`
	Formats         FileFormats  `json:"formats"`

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

// GetTokenBalanceOfOwner gets token balance of an owner
func (c *TZKT) GetTokenBalanceOfOwner(contract, tokenID, owner string) (int64, error) {
	v := url.Values{
		"account":        []string{owner},
		"token.contract": []string{contract},
		"token.tokenId":  []string{tokenID},
		"token.standard": []string{"fa2"},
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/tokens/balances/count",
		RawQuery: v.Encode(),
	}

	var balance int64

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return balance, err
	}
	if err := c.request(req, &balance); err != nil {
		return balance, err
	}

	return balance, nil
}

// GetTokenOwners returns a list of TokenOwner for a specific token
func (c *TZKT) GetTokenOwners(contract, tokenID string, limit int, lastTime time.Time) ([]TokenOwner, error) {
	v := url.Values{
		"token.contract": []string{contract},
		"token.tokenId":  []string{tokenID},
		"balance.gt":     []string{"0"},
		"token.standard": []string{"fa2"},
		"sort.asc":       []string{"lastLevel"},
		"limit":          []string{fmt.Sprintf("%d", limit)},
		"select":         []string{"account.address as address,balance,lastTime"},
	}

	rawQuery := v.Encode() + "&lastTime.ge=" + lastTime.UTC().Format(time.RFC3339)

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/tokens/balances",
		RawQuery: rawQuery,
	}

	var owners []TokenOwner

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if err := c.request(req, &owners); err != nil {
		return nil, err
	}

	return owners, nil
}

// GetTokenBalanceAndLastTimeForOwner returns balance and last activity time of an owner for a specific token
func (c *TZKT) GetTokenBalanceAndLastTimeForOwner(contract, tokenID, owner string) (int64, time.Time, error) {
	v := url.Values{
		"token.contract": []string{contract},
		"token.tokenId":  []string{tokenID},
		"balance.gt":     []string{"0"},
		"account":        []string{owner},
		"token.standard": []string{"fa2"},
		"select":         []string{"lastTime,account.address as address,balance"},
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/tokens/balances",
		RawQuery: v.Encode(),
	}

	var owners []TokenOwner

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, time.Time{}, err
	}

	if err := c.request(req, &owners); err != nil {
		return 0, time.Time{}, err
	}

	if len(owners) == 0 {
		return 0, time.Time{}, fmt.Errorf("token not found")
	}

	if len(owners) > 1 {
		return 0, time.Time{}, fmt.Errorf("multiple token owners returned")
	}

	return owners[0].Balance, owners[0].LastTime, nil
}

// GetTokenLastActivityTime returns the timestamp of the last activity for a token
func (c *TZKT) GetTokenLastActivityTime(contract, tokenID string) (time.Time, error) {
	v := url.Values{
		"token.contract": []string{contract},
		"token.tokenId":  []string{tokenID},
		"token.standard": []string{"fa2"},
		"sort.desc":      []string{"timestamp"},
		"limit":          []string{"1"},
		"select":         []string{"timestamp"},
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/tokens/transfers",
		RawQuery: v.Encode(),
	}

	var activityTime []time.Time

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return time.Time{}, err
	}

	if err := c.request(req, &activityTime); err != nil {
		return time.Time{}, err
	}

	if len(activityTime) == 0 {
		return time.Time{}, fmt.Errorf("no activities for this token")
	}

	return activityTime[0], nil
}

func (c *TZKT) GetTokenTransfers(contract, tokenID string, limit int) ([]TokenTransfer, error) {
	if limit == 0 {
		limit = 100
	}

	v := url.Values{
		"token.contract": []string{contract},
		"token.tokenId":  []string{tokenID},
		"token.standard": []string{"fa2"},
		"limit":          []string{fmt.Sprint(limit)},
		"select":         []string{"timestamp,from,to,transactionId,level"},
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/tokens/transfers",
		RawQuery: v.Encode(),
	}

	var transfers []TokenTransfer

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	if err := c.request(req, &transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

func (c *TZKT) GetTokenTransfersCount(contract, tokenID string) (int, error) {
	v := url.Values{
		"token.contract": []string{contract},
		"token.tokenId":  []string{tokenID},
		"token.standard": []string{"fa2"},
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/tokens/transfers/count",
		RawQuery: v.Encode(),
	}

	var count int

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, err
	}

	if err := c.request(req, &count); err != nil {
		return 0, err
	}

	return count, nil
}

// RetrieveTokens returns OwnedToken for a specific token. The OwnedToken object includes
// both balance and token information
func (c *TZKT) RetrieveTokens(owner string, lastTime time.Time, offset int) ([]OwnedToken, error) {
	v := url.Values{
		"account":        []string{owner},
		"limit":          []string{"50"},
		"offset":         []string{fmt.Sprintf("%d", offset)},
		"balance.ge":     []string{"0"},
		"token.standard": []string{"fa2"},
		"sort.asc":       []string{"lastLevel"},
		// NOTE: sorting over lastTime is not reliable in tzkt api. Use `lastLevel` instead
		// For example: https://api.tzkt.io/v1/tokens/balances?account=tz2GoQHhadigAa56HnAXTGAYpYn8xUZsrG11&sort=lastTime&token.standard=fa2&balance.ge=0&lastTime.ge=2022-05-16T17:09:29Z
	}

	// prevent QueryEscape for colons in time
	rawQuery := v.Encode() + "&lastTime.gt=" + lastTime.UTC().Format(time.RFC3339)

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/tokens/balances",
		RawQuery: rawQuery,
	}
	var ownedTokens []OwnedToken

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return ownedTokens, err
	}

	if err := c.request(req, &ownedTokens); err != nil {
		return ownedTokens, err
	}

	return ownedTokens, nil
}

func (c *TZKT) GetContractToken(contract, tokenID string) (Token, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.endpoint,
		Path:   "/v1/tokens",
		RawQuery: url.Values{
			"contract": []string{contract},
			"tokenId":  []string{tokenID},
		}.Encode(),
	}

	var tokenResponse []Token

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Token{}, err
	}

	if err := c.request(req, &tokenResponse); err != nil {
		return Token{}, err
	}

	if len(tokenResponse) == 0 {
		return Token{}, fmt.Errorf("token not found")
	}

	return tokenResponse[0], nil
}
