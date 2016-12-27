package bla

import "runtime"

var (
	ggPT = expvar.New("go_gc_pause_time")
)

func init() {

	ms := runtime.ReadMemStats()

}
