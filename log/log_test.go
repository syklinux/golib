package log

import (
	"testing"
)

func TestLog(t *testing.T) {

	a := LogConfig{
		Level: "INFO",
		Type:  "json",
	}

	err := InitByConf(a)
	if err != nil {
		t.Fatal(err)
	}
	Info("aaaa")
}
