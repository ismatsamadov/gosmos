package main

import "embed"

// resources holds the files gctl carries inside itself: the station table and the templates.
//
//go:embed stations.csv templates
var resources embed.FS
