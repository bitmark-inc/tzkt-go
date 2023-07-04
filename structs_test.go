package tzkt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalFlexInt64(t *testing.T) {
	var i FlexInt64

	err := json.Unmarshal([]byte(`"123"`), &i)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), int64(i))

	err = json.Unmarshal([]byte("1234"), &i)
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), int64(i))
}

func TestUnmarshalFlexBool(t *testing.T) {
	var b FlexBool

	err := json.Unmarshal([]byte(`"true"`), &b)
	assert.NoError(t, err)
	assert.Equal(t, true, bool(b))

	err = json.Unmarshal([]byte("true"), &b)
	assert.NoError(t, err)
	assert.Equal(t, true, bool(b))
}
