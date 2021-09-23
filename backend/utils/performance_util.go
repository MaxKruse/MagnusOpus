package utils

import (
	"time"

	"github.com/maxkruse/magnusopus/backend/globals"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	globals.Logger.Debugln(name, "took", elapsed)
}
