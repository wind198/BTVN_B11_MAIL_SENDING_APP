package utility

import "log"

func PrintError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
