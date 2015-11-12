package bla

import "log"

func init() {
	log.SetPrefix("[Bla] ")
}

func LFatal(arg ...interface{}) {
	if arg[0] != nil {
		log.Fatal(arg...)
	}
}

func LErr(arg ...interface{}) {

	if arg[0] != nil {
		log.Print(arg...)
	}

}

func Log(arg ...interface{}) {
	log.Print(arg...)
}
