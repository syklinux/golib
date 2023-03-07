package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

func ErrView(errno int, errmsg string, w http.ResponseWriter) {

	result := map[string]interface{}{
		"errno":  errno,
		"errmsg": errmsg,
	}

	ret, _ := json.Marshal(result)
	w.Write(ret)

	return
}

func View(result interface{}, w http.ResponseWriter, req *http.Request) {

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization,Accept,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")

	if result == nil {
		result = make(map[string]interface{}, 0)
	}
	var ret []byte
	var err error
	if res, ok := result.(map[string]interface{}); ok {

		if _, ok := res["errno"]; !ok {
			res["errno"] = 0
			res["errmsg"] = ""
		}
		if _, ok := res["server_time"]; !ok {
			res["server_time"] = int(time.Now().Unix())
		}

		ret, err = json.Marshal(res)
	} else {
		ret, err = json.Marshal(result)
	}

	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	w.Write(ret)
	return
}
