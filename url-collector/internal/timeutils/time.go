package timeutils

import "time"

const DayFormat = "2006-01-02"

// List lists all the days between given time including that days in 2006-01-01 format.
// Ex. given from 2020-01-04 to 2020-01-06 will return three dates:
// []string{"2020-01-04", "2020-01-05", "2020-01-06"}
func List(from, to time.Time) (list []string) {
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		list = append(list, d.Format(DayFormat))
	}
	return list
}
