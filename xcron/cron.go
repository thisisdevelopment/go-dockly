package xcron

import "time"

func doEvery(d time.Duration, runFunc func()) {
	runFunc()

	for range time.Tick(d) {
		runFunc()
	}
}
