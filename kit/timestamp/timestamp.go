// Package timestamp provides time/date functions.
package timestamp

import (
	"fmt"
	"math"
	"time"
)

// Age returns that of `tt`, e.g., `2 days ago`, `1 hr 36 mins ago`, etc.
func Age(tt time.Time) string {
	t := time.Duration(math.Abs(float64(time.Since(tt))))
	t = t.Round(time.Second)
	d := (t / time.Hour) / 24
	h := t / time.Hour
	t -= h * time.Hour
	m := t / time.Minute
	s := t / time.Second

	if d >= 2 {
		return fmt.Sprintf("%d days ago", d)
	}
	if d >= 1 {
		return fmt.Sprintf("%d day ago", d)
	}
	if h >= 2 {
		if m > 14 {
			return fmt.Sprintf("%d hrs %d min ago", h, m)
		}
		return fmt.Sprintf("%d hrs ago", h)
	}
	if h >= 1 {
		if m > 48 {
			return fmt.Sprintf("%d hrs ago", h+1)
		}
		if m > 9 {
			return fmt.Sprintf("%d hr %d min ago", h, m)
		}
		return fmt.Sprintf("%d hr ago", h)
	}
	if m >= 1 {
		return fmt.Sprintf("%d min ago", m)
	}
	if s >= 2 {
		return fmt.Sprintf("%d sec ago", s)
	}
	return fmt.Sprintf("1 sec ago")
}

// Unix Milliseconds @ SQL
//SELECT extract(epoch FROM NOW()) * 1000;
//SELECT (extract(epoch FROM NOW()) * 1000)::bigint;
//SELECT (extract(epoch FROM NOW())*1000)::numeric(18,0)::text;

// Truncate accuracy
// time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)

// NowUnixSec returns current Unix Time in seconds
func NowUnixSec() int64 {
	return time.Now().Unix()
}

// NowEpochMsec returns current Unix Time in milliseconds.
func NowEpochMsec() int64 {
	return time.Now().UnixNano() / 1e6
	//return time.Now().UnixNano() / int64(time.Millisecond)
	//... this other way is more commonly used, but is confusing and unnecessary.
}

// NowEpochUsec returns current Unix Time in microseconds.
func NowEpochUsec() int64 {
	return time.Now().UnixNano() / 1e3
}

// EpochSecToMsec ...
func EpochSecToMsec(t int64) int64 {
	return t * 1e3
}

// EpochMsecToSec ...
func EpochMsecToSec(t int64) int64 {
	return t / 1e3
}

// EpochSecToTimeLocal returns UTC Local
func EpochSecToTimeLocal(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// EpochMsecToTimeLocal converts millesecond Epoch to time.Time Local
func EpochMsecToTimeLocal(msec int64) time.Time {
	return time.Unix(0, msec*int64(1e6))
}

// EpochSecToTimeUTC returns GMT (Zero Offset)
func EpochSecToTimeUTC(sec int64) time.Time {
	return time.Unix(sec, 0).UTC()
}

// EpochMsecToTimeUTC converts millesecond Epoch to time.Time UTC (Zero Offset)
func EpochMsecToTimeUTC(msec int64) time.Time {
	return time.Unix(0, (msec * int64(1e6))).UTC()
} // 1614896134402 2021-03-04 22:15:34.402 +0000 UTC

// TimeStringLocal : `RFC3339` : `2020-07-22T10:21:51-04:00`
func TimeStringLocal(t time.Time) string {
	return t.Format(time.RFC3339)
}

// TimeStringZulu : `RFC3339`  : `2020-07-22T14:21:51Z`
func TimeStringZulu(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// Now returns current Unix Time in seconds.
func Now() int64 {
	return time.Now().Unix()
}

// TimeToEpochMsec ...
func TimeToEpochMsec(t time.Time) int64 {
	return t.UnixNano() / int64(1e6)
}

// TimeToEpochSec ...
func TimeToEpochSec(t time.Time) int64 {
	return t.UnixNano() / int64(1e10)
}
