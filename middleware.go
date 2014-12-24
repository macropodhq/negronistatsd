package negronistatsd

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/codegangsta/negroni"
)

// Middleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type Middleware struct {
	StatsdClient *statsd.Client
	// Name is the name of the application as recorded in latency metrics
}

// NewMiddleware returns a new *Middleware, yay!
func NewMiddleware(statsd *statsd.Client) *Middleware {
	return &Middleware{StatsdClient: statsd}
}

func (s *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	// latency in milliseconds
	latency := int64(time.Since(start) / time.Millisecond)
	res := rw.(negroni.ResponseWriter)

	s.StatsdClient.Inc("requests."+strconv.Itoa(res.Status()), 1, 1.0)
	s.StatsdClient.Timing("request_time", latency, 1.0)

}
