package errorcheck

import (
	"log"
)

func ErrorCheck(err error, message ...string) {
	if err != nil {
		if len(message) != 0 {
			// For each message we print a new line to stderr
			for _, v := range message {
				log.Println(v)
			}
			// Then write the exit code and exits
			log.Fatal(err)
		} else {
			// Write the exit code and exit
			log.Fatal(err)
		}
	}
}
