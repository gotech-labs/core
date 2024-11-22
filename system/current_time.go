package system

import (
	"time"

	"github.com/gotech-labs/core/env"
)

func Setup(tz string, offset int) func() time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.FixedZone(tz, offset)
	}
	time.Local = loc
	return time.Now
}

// CurrentTime - server current time.
var CurrentTime = Setup(env.GetString("TZ", "Asia/Tokyo"), 9*60*60)
