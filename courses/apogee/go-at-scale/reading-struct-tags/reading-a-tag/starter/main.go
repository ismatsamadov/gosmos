package main

import "fmt"

// Note is a line of text pinned beside a contact on the Board.
type Note struct {
	Text string `json:"text"`
}

func main() {
	for _, v := range []any{PipelineConfig{}, Note{}} {
		for _, line := range tagLines(v) {
			fmt.Println(line)
		}
	}
}
