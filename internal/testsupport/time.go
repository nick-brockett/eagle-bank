package testsupport

import "time"

func TimeNowRoundedMicroseconds() time.Time {
	return time.Now().Round(time.Microsecond)
}
