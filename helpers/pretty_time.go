package helpers

import (
	"fmt"
	"time"
)

func PrettyTime(t time.Time) string {

	sect := time.Now().Sub(t).Seconds()
	sec := int(sect)

	one_min := 60
	one_hour := one_min * 60
	one_day := one_hour * 24
	one_year := one_day * 365

	if sec <= one_min {
		//plural vs single
		if sec <= 1 {
			return fmt.Sprintf("%v second ago", sec)
		} else {
			return fmt.Sprintf("%v seconds ago", sec)
		}
	} else if sec > one_min && sec <= one_hour {
		//plural vs single
		min := sec / one_min
		if min <= 1 {
			return fmt.Sprintf("%v minute ago", min)
		} else {
			return fmt.Sprintf("%v minutes ago", min)
		}
	} else if sec > one_hour && sec <= one_day {
		//plural vs single
		hour := sec / one_hour
		if hour <= 1 {
			return fmt.Sprintf("%v hour ago", hour)
		} else {
			return fmt.Sprintf("%v hours ago", hour)
		}
	} else if sec > one_day && sec <= one_year {
		//plural vs single
		day := sec / one_day
		if day <= 1 {
			return fmt.Sprintf("%v day ago", day)
		} else {
			return fmt.Sprintf("%v days ago", day)
		}
	} else if sec > one_year {
		//plural vs single
		year := sec / one_year
		if year <= 1 {
			return fmt.Sprintf("%v year ago", year)
		} else {
			return fmt.Sprintf("%v years ago", year)
		}
	}

	// don't know...
	return "a while ago"
}
