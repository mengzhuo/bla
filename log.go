package bla

import "log"

func init() {
	log.SetPrefix("[Bla] ")
}

func LFatal(arg interface{}) {

	if arg != nil {
		log.Fatal(arg)
	}
}

func LErr(arg interface{}) {

	if arg != nil {
		log.Print("ERR", arg)
	}

}

func Log(arg ...interface{}) {
	log.Print(arg...)
}
