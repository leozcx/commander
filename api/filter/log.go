package filter

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

type LogFilter struct {
	requestURI string
	startTime  time.Time
}

func (filter *LogFilter) Before(context *Context, w http.ResponseWriter, r *http.Request) {
	filter.startTime = time.Now()
	filter.requestURI = r.RequestURI
}

func (filter *LogFilter) After(context *Context, w http.ResponseWriter, r *http.Request) {
	elapsedTime := time.Since(filter.startTime)
	log.WithFields(log.Fields{"method": r.Method, "uri": filter.requestURI, "elapsedTime": elapsedTime}).Info("Processing HTTP request")
}
