package tzkt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetBigMapValueByPointer returns the value of a key in a bigmap.
func (c *TZKT) GetBigMapValueByPointer(pointer int, key string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.endpoint,
		Path:   fmt.Sprintf("/v1/bigmaps/%d/keys", pointer),
		RawQuery: url.Values{
			"select": []string{"value"},
			"key":    []string{key},
		}.Encode(),
	}

	var results []json.RawMessage

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	if err := c.request(req, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("error key not found")
	}

	return results[0], nil
}

// GetBigMapPointersByContract returns a list of big map pointer for a contract.
// This call accepts tags and an option.
func (c *TZKT) GetBigMapPointersByContract(contract string, tags ...string) ([]int, error) {
	query := url.Values{
		"contract": []string{contract},
		"select":   []string{"ptr"},
	}

	if len(tags) > 0 {
		query["tags.any"] = []string{strings.Join(tags, ",")}
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.endpoint,
		Path:     "/v1/bigmaps",
		RawQuery: query.Encode(),
	}

	var pointer []int

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	if err := c.request(req, &pointer); err != nil {
		return nil, err
	}

	return pointer, nil
}

// GetBigMapsByContractAndPath get BitMap of contract
func (c *TZKT) GetBigMapsByContractAndPath(contract string, path string) (int, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.endpoint,
		Path:   "/v1/bigmaps",
		RawQuery: url.Values{
			"contract": []string{contract},
			"select":   []string{"ptr"},
			"path":     []string{path},
		}.Encode(),
	}

	var pointer []int

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, err
	}

	if err := c.request(req, &pointer); err != nil {
		return 0, err
	}

	if len(pointer) == 0 {
		return 0, fmt.Errorf("no pointer")
	}

	return pointer[0], nil
}

// GetBigMapPointerForContractTokenMetadata returns the bigmap pointer of token_metadata
// for a specific contract
func (c *TZKT) GetBigMapPointerForContractTokenMetadata(contract string) (int, error) {
	pointers, err := c.GetBigMapPointersByContract(contract, "token_metadata")
	if err != nil {
		return 0, err
	}

	if len(pointers) == 0 {
		return 0, fmt.Errorf("no pointer")
	}

	return pointers[0], nil
}
