package filter

import (
	"net/http"
)

type CorsFilter struct {
}

func writeCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

func (filter *CorsFilter) Before(context *Context, w http.ResponseWriter, r *http.Request) {
	if context.EnableCores {
		writeCorsHeaders(w, r)
	}
}

func (filter *CorsFilter) After(context *Context, w http.ResponseWriter, r *http.Request) {

}
