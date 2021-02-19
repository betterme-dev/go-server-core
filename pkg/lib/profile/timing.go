package profile

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func TimeTrack(name string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		log.Debugf("[profiling] %s took %s", name, elapsed)
	}
}
