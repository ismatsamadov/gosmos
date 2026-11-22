package main

import "net/http"

// site serves the account page under /account/ and the API at every other address, each behind its own guards.
func site(api http.Handler, a *Account, origin string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/account/", chain(a.routes(), securityHeaders, noStore, noCrossOrigin))
	mux.Handle("/", cors(origin)(api))
	return mux
}
