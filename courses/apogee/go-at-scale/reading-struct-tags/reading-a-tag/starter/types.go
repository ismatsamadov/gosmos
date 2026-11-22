package main

// PipelineConfig is tlmd's ingest settings.
type PipelineConfig struct {
	Workers int      `check:"min=1,max=32" default:"4"`
	Buffer  int      `check:"min=0,max=1024" default:"16"`
	Feeds   []string `check:"required,len=3"` // station codes
}
