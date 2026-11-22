package main

import "testing"

func TestRunKnownDecoder(t *testing.T) {
	got, err := runKnownDecoder(func() int { return 4182 })
	if got != 4182 || err != nil {
		t.Errorf("runKnownDecoder(success) = %d, %v, want 4182, nil", got, err)
	}
	got, err = runKnownDecoder(func() int { panic(decoderFault{"state lost"}) })
	if got != 0 || err == nil || err.Error() != "decoder fault: state lost" {
		t.Errorf("runKnownDecoder(known fault) = %d, %v, want 0, decoder fault: state lost", got, err)
	}
	marker := &struct{ reason string }{"unexpected"}
	for _, problem := range []any{"unrelated bug", marker, 17} {
		func() {
			defer func() {
				if got := recover(); got != problem {
					t.Errorf("runKnownDecoder(unknown) panic = %v, want original %v", got, problem)
				}
			}()
			runKnownDecoder(func() int { panic(problem) })
			t.Error("runKnownDecoder(unknown) returned, want a panic")
		}()
	}
}
