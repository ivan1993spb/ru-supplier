package main

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
)

type ErrRecurringPattern int

func (e ErrRecurringPattern) Error() string {
	return "pattern set contains recurring pattern i=" +
		strconv.Itoa(int(e))
}

type ErrInvalidPattern struct {
	err     error   // error with description
	pattern Pattern // invalid pattern
}

func (e *ErrInvalidPattern) Error() string {
	return "pattern set contains invalid pattern " +
		string(e.pattern) + "; error: " + e.err.Error()
}

type ErrFilter string

func (e ErrFilter) Error() string {
	return "filter error: " + string(e)
}

// Pattern contains regexp string
type Pattern string

func (p Pattern) Compile() (*regexp.Regexp, error) {
	return regexp.Compile(string(p))
}

// PatternSet keeps patterns for filtration
type PatternSet struct {
	prns []Pattern        // patterns
	exps []*regexp.Regexp // Regexp objects
}

func (ps PatternSet) Verify() error {
	for i := 0; i < len(ps.prns); i++ {
		// return error if there is invalid pattern
		if _, err := ps.prns[i].Compile(); err != nil {
			return &ErrInvalidPattern{err, ps.prns[i]}
		}
		for j := 0; j < len(ps.prns); j++ {
			// return error if there is recurring patterns
			if i != j && ps.prns[i] == ps.prns[j] {
				return ErrRecurringPattern(j)
			}
		}
	}
	return nil
}

// Clear removes recuring and invalid patterns
func (ps PatternSet) Clear() {
	for i := 0; i < len(ps.prns); {
		if _, err := ps.prns[i].Compile(); err != nil {
			ps.prns = append(ps.prns[:i], ps.prns[i+1:]...)
			continue
		}
		for j := 0; j < len(ps.prns); {
			if i != j && ps.prns[i] == ps.prns[j] {
				ps.prns = append(ps.prns[:j], ps.prns[j+1:]...)
			} else {
				j++
			}
		}
		i++
	}
}

func (ps PatternSet) Compile() error {
	if err := ps.Verify(); err != nil {
		return err
	}
	for _, pattern := range ps.prns {
		if exp, err := pattern.Compile(); err == nil {
			ps.exps = append(ps.exps, exp)
		} else {
			return err
		}
	}
	return nil
}

type Filter struct {
	All, OrderName, OKDP, OKPD, OrganisationName PatternSet
	fname                                        string
	enabled                                      bool
}

func LoadFilter(fname string) (*Filter, error) {
	if len(fname) == 0 {
		panic("filter: invalid file name")
	}
	file, err := os.Open(fname)
	filter := &Filter{fname: fname, enabled: true}
	if err != nil {
		if os.IsExist(err) {
			return filter, ErrFilter("can't open file: " + err.Error())
		} else {
			return filter, ErrFilter("file does not exists")
		}
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	patterns := make(map[string][]Pattern)
	err = dec.Decode(&patterns)
	if err != nil {
		if err == io.EOF {
			return filter, ErrFilter("empty JSON stream (or file)")
		} else {
			return filter, ErrFilter("invalid JSON: " + err.Error())
		}
	}
	filter.All.prns = patterns["All"]
	filter.All.Clear()
	filter.All.Compile()
	filter.OrderName.prns = patterns["OrderName"]
	filter.OrderName.Clear()
	filter.OrderName.Compile()
	filter.OKDP.prns = patterns["OKDP"]
	filter.OKDP.Clear()
	filter.OKDP.Compile()
	filter.OKPD.prns = patterns["OKPD"]
	filter.OKPD.Clear()
	filter.OKPD.Compile()
	filter.OrganisationName.prns = patterns["OrganisationName"]
	filter.OrganisationName.Clear()
	filter.OrganisationName.Compile()
	return filter, nil
}

func (f *Filter) Flush() error {
	file, err := os.Create(f.fname)
	if err != nil {
		return err
	}
	patterns := make(map[string][]Pattern)
	patterns["All"] = f.All.prns
	patterns["OrderName"] = f.OrderName.prns
	patterns["OKDP"] = f.OKDP.prns
	patterns["OKPD"] = f.OKPD.prns
	patterns["OrganisationName"] = f.OrganisationName.prns
	enc := json.NewEncoder(file)
	return enc.Encode(patterns)
}

func (f *Filter) SetEnable(flag bool) {
	f.enabled = flag
}

func (f *Filter) Exec(orders []*Order) ([]*Order, float32) {
	if !f.enabled {
		return orders, 0
	}
	count := len(orders)
	if count == 0 {
		return orders, 0
	}
	// filter all fields
	for _, exp := range f.All.exps {
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
	}
	// filter each fields
	for _, exp := range f.OrderName.exps {
		for i := 0; i < len(orders); {
			if exp.MatchString(orders[i].OrderName) {
				orders = append(orders[:i], orders[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, exp := range f.OKDP.exps {
		for i := 0; i < len(orders); {
			if exp.MatchString(orders[i].OKDP) {
				orders = append(orders[:i], orders[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, exp := range f.OKPD.exps {
		for i := 0; i < len(orders); {
			if exp.MatchString(orders[i].OKPD) {
				orders = append(orders[:i], orders[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, exp := range f.OrganisationName.exps {
		for i := 0; i < len(orders); {
			if exp.MatchString(orders[i].OrganisationName) {
				orders = append(orders[:i], orders[i+1:]...)
			} else {
				i++
			}
		}
	}
	return orders, (1 - float32(len(orders))/float32(count))
}
