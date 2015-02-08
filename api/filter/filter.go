package filter

import (
	"net/http"
)

type Filter interface {
	Before(context *Context, w http.ResponseWriter, r *http.Request)
	After(context *Context, w http.ResponseWriter, r *http.Request)
}

func NewFilter(name string) Filter {
	switch name {
	case "cors":
		return new(CorsFilter)
	case "log":
		return new(LogFilter)
	}
	return nil
}
