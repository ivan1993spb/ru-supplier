package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
)
import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
)

const (
	_STREAM_CHARSET = "windows-1251"
	_BUFFER_SIZE    = 1024
)

type ErrResponseStatus int

func (e ErrResponseStatus) Error() string {
	return "server return status " + http.StatusText(int(e))
}

type Parser struct {
	*HashStore
}

func NewParser() *Parser {
	return &Parser{LoadHashStoreSimple()}
}

func (p *Parser) Parse(resp *http.Response) ([]*Order, error) {
	if resp == nil {
		panic("parse(): passed nil response")
	}
	if resp.StatusCode != 200 {
		return nil, ErrResponseStatus(resp.StatusCode)
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
	checking_chunk, exists := p.GetHashChunk(rawurl)
	// current newest chunk will checking chunk at next time
	newest_chunk, err := brdr.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// cut delim \n
	newest_chunk = newest_chunk[:len(newest_chunk)-1]
	// check for updates if exists data in hashstore
	if exists {
		hash := md5.New()
		hash.Write(newest_chunk)
		if bytes.Compare(checking_chunk, hash.Sum(nil)) == 0 {
			// return if feed was not updated
			return nil, nil
		}
	}
	// save newest chunk in cache
	p.SetHashChunk(rawurl, newest_chunk)
	if err = p.Save(); err != nil {
		log.Println("can't save cache:", err)
	}
	// convert newst chunk to order
	row := regexp.MustCompile(`\s*;\s*`).
		Split(string(newest_chunk), _ORDER_COLUMN_COUNT)
	if len(row) != _ORDER_COLUMN_COUNT {
		return nil, errors.New("invalud column count")
	}
	// save to order list
	var order_row [_ORDER_COLUMN_COUNT]string
	copy(order_row[:], row)
	orders := []*Order{NewOrder(order_row)}
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
	rdr.FieldsPerRecord = _ORDER_COLUMN_COUNT
	for {
		row, err = rdr.Read()
		if err != nil && err != io.EOF {
			return orders, err
		}
		if err == io.EOF && len(row) == 0 {
			break
		}
		copy(order_row[:], row)
		// csv reader checks count of fields
		orders = append(orders, NewOrder(order_row))
		if err == io.EOF {
			break
		}
	}
	return orders, nil
}

func (p *Parser) RemoveCache() error {
	return p.Remove()
}
