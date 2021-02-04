package common

import "log"

//PanicOnError if err != nil log.Panic
func PanicOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
