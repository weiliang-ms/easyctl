package format

import (
	"bytes"
	"encoding/json"
)

func Object(v interface{}) (bytes.Buffer, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return bytes.Buffer{}, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, bs, "", "\t")
	if err != nil {
		return bytes.Buffer{}, err
	}

	return out, nil
}
