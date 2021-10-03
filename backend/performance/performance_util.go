package performance

import (
	"log"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Println(name, "took", elapsed)
}
