package utils

import (
	"context"
	"net/http"
	"time"
)

const ContextKey = "request_context"

type Context struct {
	Stores map[string]interface{}
}

func NewContext() *Context {
	ctx := &Context{}
	ctx.Stores = make(map[string]interface{})

	return ctx
}

func SetReqContext(req *http.Request, ctx *Context) *http.Request {
	httpContext := context.WithValue(req.Context(), ContextKey, ctx)
	return req.WithContext(httpContext)
}

func GetReqContext(r *http.Request) *Context {
	v := r.Context().Value(ContextKey)
	if v == nil {
		return nil
	}

	return v.(*Context)
}

func (c *Context) Set(key string, value interface{}) {
	if c.Stores == nil {
		c.Stores = make(map[string]interface{})
	}
	c.Stores[key] = value
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Stores[key]
	return
}

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}
