package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
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

func Parse(resp *http.Response) ([]*Order, error) {
	if resp == nil {
		panic("parse(): passed nil response")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("server return status:" + resp.Status)
	}
	w1251rdr, err := charset.NewReader(_STREAM_CHARSET, resp.Body)
	if err != nil {
		return nil, err
	}
	brdr := bufio.NewReaderSize(w1251rdr, _BUFFER_SIZE)
	// skip first line with topics
	if _, err = brdr.ReadString('\n'); err != nil {
		return nil, errors.New("skip first line err: " + err.Error())
	}
	// below get checking chunk by url
	// get url
	rawurl := resp.Request.URL.String()
	// get checking chunk
	// checking chunk was newest chunk at last time
	checking_chunk, exists := hashstore.GetHashChunk(rawurl)
	// current newest chunk will checkin chunk at next time
	newest_chunk, err := brdr.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// cut delim \n
	newest_chunk = newest_chunk[:len(newest_chunk)-1]
	// return if feed was not updated
	nhash := md5.New()
	nhash.Write(newest_chunk)
	if bytes.Compare(checking_chunk, nhash.Sum(nil)) == 0 {
		return nil, nil
	}
	// save newest chunk in cache
	hashstore.SetHashChunk(rawurl, newest_chunk)
	if err = hashstore.Flush(); err != nil {
		log.Error.Println("hashstore error:", err)
	}
	// convert newst chunk to order
	order_row := regexp.MustCompile(`\s*;\s*`).
		Split(string(newest_chunk), _COLUMN_COUNT)
	if len(order_row) != _COLUMN_COUNT {
		return nil, errors.New("invalud column count")
	}
	// save to order list
	orders := []*Order{NewOrder(order_row[0], order_row[1],
		order_row[2], order_row[3], order_row[4], order_row[5],
		order_row[6], order_row[7], order_row[8], order_row[9],
		order_row[10], order_row[11], order_row[12], order_row[13],
		order_row[14], order_row[15], order_row[16])}
	// setting csv reader
	var rdr *csv.Reader
	// if exists checking chunk read while does not find matched chunk
	if exists && len(checking_chunk) > 0 {
		rdr = csv.NewReader(NewCacheReader(
			brdr, []byte{'\n'}, checking_chunk,
		))
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
		// csv reader checks count of fields
		orders = append(orders, NewOrder(row[0], row[1], row[2], row[3],
			row[4], row[5], row[6], row[7], row[8], row[9], row[10],
			row[11], row[12], row[13], row[14], row[15], row[16]))
		if err == io.EOF {
			break
		}
	}
	return orders, nil
}
