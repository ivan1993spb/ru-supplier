package main

import (
	"encoding/csv"
	"io"
)
import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
)

const (
	_COLUMN_COUNT   = 17
	_STREAM_CHARSET = "windows-1251"
)

func parse(r io.Reader) ([]*Order, error) {
	w1251rdr, err := charset.NewReader(_STREAM_CHARSET, r)
	if err != nil {
		return nil, err
	}
	rdr := csv.NewReader(w1251rdr)
	rdr.Comma = ';'
	rdr.TrimLeadingSpace = true
	rdr.TrailingComma = true
	rdr.FieldsPerRecord = _COLUMN_COUNT
	// skip first line
	rdr.Read()
	orders := make([]*Order, 0)
	for {
		row, err := rdr.Read()
		if err != nil && err != io.EOF {
			return orders, err
		}
		if err == io.EOF && len(row) == 0 {
			break
		}
		orders = append(orders, NewOrder(row[0], row[1],
			row[2], row[3], row[4], row[5], row[6], row[7],
			row[8], row[9], row[10], row[11], row[12], row[13],
			row[14], row[15], row[16]))
		if err == io.EOF {
			break
		}
	}
	return orders, nil
}
