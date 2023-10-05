package tzkt

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GetLevelByTime returns the block level of the given time
func (c *TZKT) GetLevelByTime(at time.Time) (uint64, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.endpoint,
		Path:   fmt.Sprintf("/v1/blocks/%s/level", at.UTC().Format(time.RFC3339)),
	}

	var level uint64

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, err
	}
	if err := c.request(req, &level); err != nil {
		return 0, err
	}

	return level, nil
}
