package main

import "runtime"

// inUse returns how many bytes of heap the program still uses, counted after a garbage collection.
func inUse() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}
