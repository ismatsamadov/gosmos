// The courses' text holds every lesson's Go files, which build in the app's lesson folders,
// never here: this go.mod keeps them out of the repository's module, so go vet ./... and
// go test ./... pass them by.
module courses
