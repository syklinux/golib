package mysql

import (
	"context"
	"strconv"

	"github.com/nobugtodebug/go-objectid"
	"github.com/smallnest/rpcx/server"
)

func AttachTraceId(c context.Context) string {
	if c == nil {
		c = context.Background()
	}

	var traceId string
	var cv = c.Value(server.StartRequestContextKey)
	if value, ok := cv.(int64); ok {
		traceId = strconv.FormatInt(value, 10)
	}

	if traceId == "" {
		traceId = objectid.New().String()
	}

	return traceId
}
