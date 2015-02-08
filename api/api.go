package api

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/denverdino/commander/api/filter"
	"github.com/denverdino/commander/registry"
	"github.com/gorilla/mux"
	"net/http"
)

const APIVERSION = "1.16"

//type context struct {
//	addr      string
//	debug     bool
//	version   string
//	tlsConfig *tls.Config
//}

//type handler func(c *context, w http.ResponseWriter, r *http.Request)

// Default handler for methods not supported by clustering.
func notImplementedHandler(c *filter.Context, w http.ResponseWriter, r *http.Request) int {
	status := http.StatusNotImplemented
	httpError(w, "Not supported in clustering mode.", status)
	return status
}

func httpError(w http.ResponseWriter, err string, status int) {
	log.WithField("status", status).Errorf("HTTP error: %v", err)
	http.Error(w, err, status)
}

// Proxy a request to Docker Swarm
func proxyRequest(c *filter.Context, w http.ResponseWriter, r *http.Request) int {
	status, err := proxy(c.TLSConfig, c.Addr, w, r)
	if err != nil {

		httpError(w, err.Error(), status)
	}
	return status
}

func createRouter(c *filter.Context) *mux.Router {
	r := mux.NewRouter()

	handleServiceResources(c, "/services/create", "POST", r)

	r.PathPrefix("/").HandlerFunc(wrapFunc(c, proxyRequest))
	//r.HandleFunc(wrapFunc(c, proxyRequest))
	return r
}

func wrapFunc(c *filter.Context, f HTTPHandlerFunc) http.HandlerFunc {
	wrap := func(w http.ResponseWriter, r *http.Request) {
		f(c, w, r)
	}
	return wrap
}

func createServiceResource(c *filter.Context, w http.ResponseWriter, r *http.Request) int {

	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err == nil {
		serviceId := data["id"].(string)
		fmt.Printf("data = %v\n", data)

		registry.SetService(c.EtcdClient, serviceId, &data)
	} else {
		log.Error(err)
	}

	return http.StatusCreated
}

func handleServiceResources(c *filter.Context, localRoute string, localMethod string, router *mux.Router) {
	wrap := wrapFunc(c, createServiceResource)
	router.Path("/v{version:[0-9.]+}" + localRoute).Methods(localMethod).HandlerFunc(wrap)
	router.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)
}
