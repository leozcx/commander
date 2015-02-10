package filter

import (
	"github.com/denverdino/commander/context"
	"net/http"
)

type Filter interface {
	Before(context *context.Context, w http.ResponseWriter, r *http.Request)
	After(context *context.Context, w http.ResponseWriter, r *http.Request)
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
