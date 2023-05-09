package tzkt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type NullableInt int64

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

// StringBool is a type that accepts a bool can be either string or boolean
type StringBool bool

func (b *StringBool) UnmarshalJSON(data []byte) error {
	var _b bool

	err := json.Unmarshal(bytes.Trim(data, `"`), &_b)
	if err != nil {
		return nil
	}

	*b = StringBool(_b)

	return nil
}

type MIMEFormat string

func (m *MIMEFormat) UnmarshalJSON(data []byte) error {
	if data[0] == 91 {
		data = bytes.Trim(data, "[]")
	}

	return json.Unmarshal(data, (*string)(m))
}

type FileFormats []FileFormat

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
