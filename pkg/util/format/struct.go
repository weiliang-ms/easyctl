package format

import (
	"bytes"
	"encoding/json"
)

func ObjectToJson(v interface{}) bytes.Buffer {
	bs, _ := json.Marshal(v)
	var out bytes.Buffer
	_ = json.Indent(&out, bs, "", "\t")

	return out
}
