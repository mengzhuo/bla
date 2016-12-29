package bla

import "log"

func ErrOrOk(action string, err error) {
	if err != nil {
		log.Printf("%s ... [ERR] reason:%v", action, err)
	} else {
		log.Print(action, " ... [OK]")
	}
}
