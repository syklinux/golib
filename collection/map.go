package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

func Map2Buffer(mm map[string]interface{}) *bytes.Buffer {
	data, err := json.Marshal(mm)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer([]byte(data))
}

func Map2UrlParams(mm map[string]interface{}) string {
	vals := url.Values{}
	for key, val := range mm {
		vals.Add(key, fmt.Sprintf("%v", val))
	}
	return vals.Encode()
}
