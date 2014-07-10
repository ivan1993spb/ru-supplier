package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"regexp"
)
import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
)

const (
	_COLUMN_COUNT   = 17
	_STREAM_CHARSET = "windows-1251"
	_BUFFER_SIZE    = 1024
)

func parse(resp *http.Response) ([]*Order, error) {
	w1251rdr, err := charset.NewReader(_STREAM_CHARSET, resp.Body)
	if err != nil {
		return nil, err
	}
	brdr := bufio.NewReaderSize(w1251rdr, _BUFFER_SIZE)
	// skip first topic line
	brdr.ReadString('\n')

	rawurl := resp.Request.URL.String()
	// newest chunk at last time
	checking_chunk, exists := hashstore.GetHashChunk(rawurl)
	// newest chunk now
	newest_chunk, err := brdr.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// delete delim \n
	newest_chunk = newest_chunk[:len(newest_chunk)-1]
	// save newest chunk in cache
	hashstore.SetHashChunk(rawurl, newest_chunk)
	if err = hashstore.Flush(); err != nil {
		log.Error.Println("hashstore error:", err)
	}
	// convert newst chunk to order
	exp := regexp.MustCompile(`\s*;\s*`)
	order_row := exp.Split(string(newest_chunk), _COLUMN_COUNT)
	if len(order_row) != _COLUMN_COUNT {
		return nil, errors.New("Invalud column count")
	}
	// save to order list
	orders := []*Order{NewOrder(order_row[0], order_row[1],
		order_row[2], order_row[3], order_row[4], order_row[5],
		order_row[6], order_row[7], order_row[8], order_row[9],
		order_row[10], order_row[11], order_row[12], order_row[13],
		order_row[14], order_row[15], order_row[16])}
	// setting of reader
	var rdr *csv.Reader
	// if exists checking chunk read while does not find matched chunk
	if exists && len(checking_chunk) > 0 {
		rdr = csv.NewReader(
			NewCacheReader(brdr, []byte{'\n'}, checking_chunk),
		)
	} else {
		rdr = csv.NewReader(brdr)
	}
	rdr.Comma = ';'
	rdr.TrimLeadingSpace = true
	rdr.TrailingComma = true
	rdr.FieldsPerRecord = _COLUMN_COUNT
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
