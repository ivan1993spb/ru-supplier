package main

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
)

const _FILE_FILTERS = "filters.json"

type ErrInvalidRegexp struct {
	sect, ptrn string // section and pattern
	err        error
}

func (e *ErrInvalidRegexp) Error() string {
	return `invalid regexp filter in "` + e.sect + `": ` +
		e.err.Error() + `; ` + e.ptrn
}

func filter(orders []*Order) []*Order {
	if len(orders) == 0 {
		return orders
	}
	file, err := os.Open(_FILE_FILTERS)
	if err != nil {
		if os.IsExist(err) {
			log.Error.Println("filter(): can't open file with",
				"filter patterns:", err)
		} else {
			log.Warning.Println(_FILE_FILTERS, "was not found")
		}
		return orders
	}
	dec := json.NewDecoder(file)
	var patterns *struct {
		All, OrderName, OKDP, OKPD, OrganisationName []string
	}
	err = dec.Decode(&patterns)
	if err != nil {
		if err != io.EOF {
			log.Error.Println("filter(): can't decode json",
				"from file with filter patterns:", err)
		} else {
			log.Warning.Println("empty filter file", _FILE_FILTERS)
		}
		return orders
	}
	// filter all fields
	for _, pattern := range patterns.All {
		exp, err := regexp.Compile(pattern)
		if err == nil {
			for i := 0; i < len(orders); {
				if exp.MatchString(orders[i].OrderName) ||
					exp.MatchString(orders[i].OKDP) ||
					exp.MatchString(orders[i].OKPD) ||
					exp.MatchString(orders[i].OrganisationName) {
					orders = append(orders[:i], orders[i+1:]...)
				} else {
					i++
				}
			}
		} else {
			log.Warning.Println(&ErrInvalidRegexp{
				"All",
				pattern,
				err,
			})
		}
	}
	// filter each fields
	for _, pattern := range patterns.OrderName {
		exp, err := regexp.Compile(pattern)
		if err == nil {
			for i := 0; i < len(orders); {
				if exp.MatchString(orders[i].OrderName) {
					orders = append(orders[:i], orders[i+1:]...)
				} else {
					i++
				}
			}
		} else {
			log.Warning.Println(&ErrInvalidRegexp{
				"OrderName",
				pattern,
				err,
			})
		}
	}
	for _, pattern := range patterns.OKDP {
		exp, err := regexp.Compile(pattern)
		if err == nil {
			for i := 0; i < len(orders); {
				if exp.MatchString(orders[i].OKDP) {
					orders = append(orders[:i], orders[i+1:]...)
				} else {
					i++
				}
			}
		} else {
			log.Warning.Println(&ErrInvalidRegexp{
				"OKDP",
				pattern,
				err,
			})
		}
	}
	for _, pattern := range patterns.OKPD {
		exp, err := regexp.Compile(pattern)
		if err == nil {
			for i := 0; i < len(orders); {
				if exp.MatchString(orders[i].OKPD) {
					orders = append(orders[:i], orders[i+1:]...)
				} else {
					i++
				}
			}
		} else {
			log.Warning.Println(&ErrInvalidRegexp{
				"OKPD",
				pattern,
				err,
			})
		}
	}
	for _, pattern := range patterns.OrganisationName {
		exp, err := regexp.Compile(pattern)
		if err == nil {
			for i := 0; i < len(orders); {
				if exp.MatchString(orders[i].OrganisationName) {
					orders = append(orders[:i], orders[i+1:]...)
				} else {
					i++
				}
			}
		} else {
			log.Warning.Println(&ErrInvalidRegexp{
				"OrganisationName",
				pattern,
				err,
			})
		}
	}
	return orders
}
