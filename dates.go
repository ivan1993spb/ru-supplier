package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// zakupki.goc.ru uses moscow time zone
var MoscowTimeZone = time.FixedZone("MSK", 4*60*60)

func ParseRusFormatDate(date string) (time.Time, error) {
	var null time.Time
	if len(date) > 0 {
		chunks := strings.SplitN(date, ".", 3)
		if len(chunks) >= 3 {
			if day, err := strconv.Atoi(chunks[0]); err == nil {
				if month, err := strconv.Atoi(chunks[1]); err == nil {
					if year, err := strconv.Atoi(chunks[2]); err == nil {
						return time.Date(year, time.Month(month), day,
							0, 0, 0, 0, MoscowTimeZone), nil
					}
				}
			}
		}
	}
	return null, errors.New("Invalid russian date format")
}

func RusFormatDate(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%0.2d.%0.2d.%0.4d", d, m, y)
}
