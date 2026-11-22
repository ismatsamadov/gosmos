package main

import "net/http"

// stationsHandler returns the handler that answers with the service's stations.
func (s *Service) stationsHandler() http.Handler {
	return http.NotFoundHandler()
}
