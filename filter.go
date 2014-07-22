package main

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
)

type ErrInvalidPattern struct {
	err     error   // error with description
	pattern Pattern // invalid pattern
}

func (e *ErrInvalidPattern) Error() string {
	return "Invalid pattern " + string(e.pattern) + ": " + e.err.Error()
}

type ErrRecurringPattern Pattern

func (e ErrRecurringPattern) Error() string {
	return "Recurring pattern " + string(e)
}

type Pattern string

func (p Pattern) Compile() (*regexp.Regexp, error) {
	return regexp.Compile(string(p))
}

type PatternSet []Pattern

func (ps PatternSet) Clear() {
	for len(ps) > 0 {
		if i, _ := ps.Verify(); i > -1 {
			ps = append(ps[:i], ps[i+1:]...)
		} else {
			return
		}
	}
}

// Verify returns error and pattern index if there is an error
func (ps PatternSet) Verify() (int, error) {
	for i := 0; i < len(ps); i++ {
		// return error if there is invalid pattern
		if _, err := ps[i].Compile(); err != nil {
			return i, &ErrInvalidPattern{err, ps[i]}
		}
		for j := 0; j < len(ps); j++ {
			// return error if there is recurring patterns
			if i != j && ps[i] == ps[j] {
				return j, ErrRecurringPattern(ps[j])
			}
		}
	}
	return -1, nil
}

func (ps PatternSet) Compile() (ExpSet, error) {
	es := make([]*regexp.Regexp, len(ps))
	for i, pattern := range ps {
		if exp, err := pattern.Compile(); err == nil {
			es[i] = exp
		} else {
			return nil, err
		}
	}
	return es, nil
}

type ExpSet []*regexp.Regexp

func (es ExpSet) PatternSet() PatternSet {
	ps := make(PatternSet, len(es))
	for i, exp := range es {
		ps[i] = Pattern(exp.String())
	}
	return ps
}

type Filter struct {
	All, OrderName, OKDP, OKPD, OrganisationName ExpSet
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
	patterns := make(map[string]PatternSet)
	err = dec.Decode(&patterns)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return
	}
	if _, ok := patterns["All"]; ok && len(patterns["All"]) > 0 {
		filter.SetExpsAll(patterns["All"])
	}
	if _, ok := patterns["OrderName"]; ok && len(patterns["OrderName"]) > 0 {
		filter.SetExpsOrderName(patterns["OrderName"])
	}
	if _, ok := patterns["OKDP"]; ok && len(patterns["OKDP"]) > 0 {
		filter.SetExpsOKDP(patterns["OKDP"])
	}
	if _, ok := patterns["OKPD"]; ok && len(patterns["OKPD"]) > 0 {
		filter.SetExpsOKPD(patterns["OKPD"])
	}
	if _, ok := patterns["OrganisationName"]; ok && len(patterns["OrganisationName"]) > 0 {
		filter.SetExpsOrganisationName(patterns["OrganisationName"])
	}
	return filter, nil
}

func (f *Filter) SetExpsAll(ps PatternSet) {
	if len(ps) > 0 {
		ps.Clear()
		if len(ps) > 0 {
			f.All, _ = ps.Compile()
		}
	}
}

func (f *Filter) SetExpsOrderName(ps PatternSet) {
	if len(ps) > 0 {
		ps.Clear()
		if len(ps) > 0 {
			f.OrderName, _ = ps.Compile()
		}
	}
}

func (f *Filter) SetExpsOKDP(ps PatternSet) {
	if len(ps) > 0 {
		ps.Clear()
		if len(ps) > 0 {
			f.OKDP, _ = ps.Compile()
		}
	}
}

func (f *Filter) SetExpsOKPD(ps PatternSet) {
	if len(ps) > 0 {
		ps.Clear()
		if len(ps) > 0 {
			f.OKPD, _ = ps.Compile()
		}
	}
}

func (f *Filter) SetExpsOrganisationName(ps PatternSet) {
	if len(ps) > 0 {
		ps.Clear()
		if len(ps) > 0 {
			f.OrganisationName, _ = ps.Compile()
		}
	}
}

// Execute executes filter for order list and returns statistic
func (f *Filter) Execute(orders []*Order) ([]*Order, float32) {
	count := len(orders)
	if count == 0 {
		return orders, 0
	}
	// filter all fields
	for _, exp := range f.All {
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
	for _, exp := range f.OrderName {
		for i := 0; i < len(orders); {
			if exp.MatchString(orders[i].OrderName) {
				orders = append(orders[:i], orders[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, exp := range f.OKDP {
		for i := 0; i < len(orders); {
			if exp.MatchString(orders[i].OKDP) {
				orders = append(orders[:i], orders[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, exp := range f.OKPD {
		for i := 0; i < len(orders); {
			if exp.MatchString(orders[i].OKPD) {
				orders = append(orders[:i], orders[i+1:]...)
			} else {
				i++
			}
		}
	}
	for _, exp := range f.OrganisationName {
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

func (f *Filter) Save() error {
	file, err := os.Create(f.fname)
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(map[string]PatternSet{
		"All":              f.All.PatternSet(),
		"OrderName":        f.OrderName.PatternSet(),
		"OKDP":             f.OKDP.PatternSet(),
		"OKPD":             f.OKPD.PatternSet(),
		"OrganisationName": f.OrganisationName.PatternSet(),
	})
}
