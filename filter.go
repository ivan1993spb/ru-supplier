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
}

func LoadFilter(fname string) (filter *Filter, err error) {
	if len(fname) == 0 {
		panic("filter: invalid file name")
	}
	var file *os.File
	file, err = os.Open(fname)
	filter = &Filter{fname: fname}
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	patterns := make(map[string][]Pattern)
	err = dec.Decode(&patterns)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return
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

func (f *Filter) Save() error {
	file, err := os.Create(f.fname)
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(map[string][]Pattern{
		"All":              f.All.prns,
		"OrderName":        f.OrderName.prns,
		"OKDP":             f.OKDP.prns,
		"OKPD":             f.OKPD.prns,
		"OrganisationName": f.OrganisationName.prns,
	})
}

// Execute executes filter for order list and returns statistic
func (f *Filter) Execute(orders []*Order) ([]*Order, float32) {
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
