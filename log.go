package bla

import "log"

func init() {
	log.SetPrefix("[Bla] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
