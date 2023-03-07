package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
)

func GetAllParams(req *http.Request) (ps map[string]string, err error) {
	err = req.ParseForm()
	if err != nil {
		return nil, err
	}

	ps = map[string]string{}
	for name, _ := range req.Form {
		ps[name] = req.Form.Get(name)
	}

	ct := req.Header.Get("Content-Type")
	if ct != "" {
		ct, _, err = mime.ParseMediaType(ct)
		if ct == "multipart/form-data" {
			err = req.ParseMultipartForm(4096)
			if err != nil {
				return nil, err
			}
			for name, _ := range req.Form {
				len := len(req.Form[name])
				if len == 0 {
					ps[name] = ""
				} else {
					ps[name] = req.Form[name][len-1]
				}
			}
		} else if ct == "application/json" {
			var body map[string]interface{}
			bs, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewReader(bs))
			if bs != nil && len(bs) != 0 {
				if err = json.Unmarshal(bs, &body); err != nil {
					return
				}
				for k, v := range body {
					ps[k] = ToString(v)
				}
			}
		}
	}

	return ps, nil
}
