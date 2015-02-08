package api

import (
	"github.com/denverdino/commander/api/filter"
	"net/http"
)

type HTTPHandlerFunc func(c *filter.Context, w http.ResponseWriter, r *http.Request) int

type Interceptor struct {
	context *filter.Context
	filters []filter.Filter
	handler http.Handler
}

func NewInterceptor(context *filter.Context, handler http.Handler) *Interceptor {
	return &Interceptor{
		context: context,
		handler: handler,
	}
}

// Add a filter to Interceptor
func (interceptor *Interceptor) addFilter(filter filter.Filter) *Interceptor {
	interceptor.filters = append(interceptor.filters, filter)
	return interceptor
}

func (interceptor *Interceptor) addFilterByName(name string) *Interceptor {
	filter := filter.NewFilter(name)
	interceptor.filters = append(interceptor.filters, filter)
	return interceptor
}

func (interceptor *Interceptor) GetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, filter := range interceptor.filters {
			filter.Before(interceptor.context, w, r)
		}
		interceptor.handler.ServeHTTP(w, r)
		for i := len(interceptor.filters) - 1; i >= 0; i-- {
			filter := interceptor.filters[i]
			filter.After(interceptor.context, w, r)
		}
	})
}
