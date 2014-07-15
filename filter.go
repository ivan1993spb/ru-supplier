package main

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
)

// type ErrInvalidRegexp struct {
// 	sect, ptrn string // section and pattern
// 	err        error
// }

// func (e *ErrInvalidRegexp) Error() string {
// 	return `invalid regexp filter in "` + e.sect + `": ` +
// 		e.err.Error() + `; ` + e.ptrn
// }

// Pattern contains regexp string
type Pattern string

func (p *Pattern) Compile(pattern) (*regexp.Regexp, error) {
	return regexp.Compile(string(p))
}

type PatternSet []Pattern

func (ps PatternSet) Add(p Pattern) bool {
	for _, pattern := range ps {
		if pattern == p {
			return false
		}
	}
	ps = append(ps, p)
	return true
}

func (ps PatternSet) Del(i int) bool {
	if len(ps) > i {
		return false
	}
	ps = append(ps[:i], ps[i+1:]...)
	return true
}

// PatternSet keeps patterns for filtration
type Filter struct {
	All, OrderName, OKDP, OKPD, OrganisationName PatternSet
	fname                                        string
	enabled                                      bool
}

func ReadFilter(r io.Reader) (f *Filter, err error) {
	dec := json.NewDecoder(r)
	err = dec.Decode(&f)
	return
}

func (f *Filter) Flush() error {
	file, err := os.Create(f.fname)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	return enc.Encode(f)
}

// Verify checks patterns. Return error if is there or nil
// func (ps *PatternSet) Verify() error {

// }

// type Filter struct {
// 	fname   string
// 	enabled bool
// }

// func LoadFilter(fname string) (*Filter, error) {
// 	if len(fname) == 0 {
// 		panic("filter: invalid file name")
// 	}
// 	file, err := os.Open(fname)
// 	if err != nil {
// 		return
// 	}
// 	return &Filter{fname, true}
// }

// func (f *Filter) SetEnable(flag bool) {
// 	f.enabled = flag
// }

// func (f *Filter) Filter(orders []*Order) []*Order {
// 	if !f.enabled {
// 		return orders
// 	}
// 	count := len(orders)
// 	if count == 0 {
// 		return orders
// 	}
// 	file, err := os.Open(f.fname)
// 	if err != nil {
// 		if os.IsExist(err) {
// 			log.Error.Println("filter(): can't open file with",
// 				"filter patterns:", err)
// 		} else {
// 			log.Warning.Println(f.fname, "was not found")
// 		}
// 		return orders
// 	}
// 	dec := json.NewDecoder(file)
// 	var patterns *struct {
// 		All, OrderName, OKDP, OKPD, OrganisationName []string
// 	}
// 	err = dec.Decode(&patterns)
// 	if err != nil {
// 		if err != io.EOF {
// 			log.Error.Println("filter(): can't decode json",
// 				"from file with filter patterns:", err)
// 		} else {
// 			log.Warning.Println("empty filter file", f.fname)
// 		}
// 		return orders
// 	}
// 	// filter all fields
// 	for _, pattern := range patterns.All {
// 		exp, err := regexp.Compile(pattern)
// 		if err == nil {
// 			for i := 0; i < len(orders); {
// 				if exp.MatchString(orders[i].OrderName) ||
// 					exp.MatchString(orders[i].OKDP) ||
// 					exp.MatchString(orders[i].OKPD) ||
// 					exp.MatchString(orders[i].OrganisationName) {
// 					orders = append(orders[:i], orders[i+1:]...)
// 				} else {
// 					i++
// 				}
// 			}
// 		} else {
// 			log.Warning.Println(&ErrInvalidRegexp{
// 				"All",
// 				pattern,
// 				err,
// 			})
// 		}
// 	}
// 	// filter each fields
// 	for _, pattern := range patterns.OrderName {
// 		exp, err := regexp.Compile(pattern)
// 		if err == nil {
// 			for i := 0; i < len(orders); {
// 				if exp.MatchString(orders[i].OrderName) {
// 					orders = append(orders[:i], orders[i+1:]...)
// 				} else {
// 					i++
// 				}
// 			}
// 		} else {
// 			log.Warning.Println(&ErrInvalidRegexp{
// 				"OrderName",
// 				pattern,
// 				err,
// 			})
// 		}
// 	}
// 	for _, pattern := range patterns.OKDP {
// 		exp, err := regexp.Compile(pattern)
// 		if err == nil {
// 			for i := 0; i < len(orders); {
// 				if exp.MatchString(orders[i].OKDP) {
// 					orders = append(orders[:i], orders[i+1:]...)
// 				} else {
// 					i++
// 				}
// 			}
// 		} else {
// 			log.Warning.Println(&ErrInvalidRegexp{
// 				"OKDP",
// 				pattern,
// 				err,
// 			})
// 		}
// 	}
// 	for _, pattern := range patterns.OKPD {
// 		exp, err := regexp.Compile(pattern)
// 		if err == nil {
// 			for i := 0; i < len(orders); {
// 				if exp.MatchString(orders[i].OKPD) {
// 					orders = append(orders[:i], orders[i+1:]...)
// 				} else {
// 					i++
// 				}
// 			}
// 		} else {
// 			log.Warning.Println(&ErrInvalidRegexp{
// 				"OKPD",
// 				pattern,
// 				err,
// 			})
// 		}
// 	}
// 	for _, pattern := range patterns.OrganisationName {
// 		exp, err := regexp.Compile(pattern)
// 		if err == nil {
// 			for i := 0; i < len(orders); {
// 				if exp.MatchString(orders[i].OrganisationName) {
// 					orders = append(orders[:i], orders[i+1:]...)
// 				} else {
// 					i++
// 				}
// 			}
// 		} else {
// 			log.Warning.Println(&ErrInvalidRegexp{
// 				"OrganisationName",
// 				pattern,
// 				err,
// 			})
// 		}
// 	}
// 	filtered := (1 - float32(len(orders))/float32(count))
// 	log.Warning.Printf("filtered %.1f%%\n", filtered*100)
// 	return orders
// }
