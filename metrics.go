package metrics

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New() http.Handler {
	var client *gophercloud.ServiceClient

	r := prometheus.NewRegistry()
	r.MustRegister(
		baremetalService(client),
		blockstorageService(client),
		computeService(client),
	)

	return promhttp.HandlerFor(r, promhttp.HandlerOpts{})
}
