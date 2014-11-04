package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"log"
	"net/http"

	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
)

const (
	_STREAM_CHARSET = "windows-1251"
	_BUFFER_SIZE    = 1536
)

type OrderParserReader interface {
	ReadOrders(*http.Response) ([]*Order, error)
	RemoveCache() error
}

type OrderReader struct {
	*HashStore
}

func NewOrderReader() *OrderReader {
	return &OrderReader{LoadHashStoreSimple()}
}

func (p *OrderReader) ReadOrders(resp *http.Response) (
	[]*Order, error) {

	if resp == nil {
		panic("ReadOrders(): passed nil response")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Server return status " + resp.Status)
	}

	w1251rdr, err := charset.NewReader(_STREAM_CHARSET, resp.Body)
	if err != nil {
		return nil, err
	}
	brdr := bufio.NewReaderSize(w1251rdr, _BUFFER_SIZE)

	// skip first line with topics
	if _, err = brdr.ReadString('\n'); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, errors.New("Skip first line err: " + err.Error())
	}

	// Below get newest chunk from stream and checking chunk
	// from store. Current newest chunk will checking chunk at next
	// time
	newestChunk, err := brdr.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}
	// There is no orders. Empty result
	if len(newestChunk) == 0 && err == io.EOF {
		return nil, nil
	}
	if newestChunk[len(newestChunk)-1] == '\n' {
		// cut delim \n if there is
		newestChunk = newestChunk[:len(newestChunk)-1]
	}
	// below get checking chunk by url
	rawurl := resp.Request.URL.String()
	// checking chunk was newest chunk at last time
	checkingChunk, exists := p.HashStore.GetHashChunk(rawurl)
	// check for updates if exists data in hashstore by comparing
	// newest chunk and checking chunk
	if exists {
		hash := md5.New()
		hash.Write(newestChunk)
		if bytes.Compare(checkingChunk, hash.Sum(nil)) == 0 {
			// return if feed was not updated
			return nil, nil
		}
	}
	// save newest chunk in cache
	p.HashStore.SetHashChunk(rawurl, newestChunk)
	if err = p.HashStore.Save(); err != nil {
		log.Println("Can't save cache:", err)
	}

	var orders []*Order
	if order, err := ParseOrder(newestChunk); err == nil {
		orders = append(orders, order)
	} else {
		log.Println("Parsing order error:", err)
	}

	// if exists checking chunk read while does not find matched chunk
	if exists {
		brdr = bufio.NewReader(NewCacheReader(
			brdr, []byte{'\n'}, checkingChunk,
		))
	}

	for {
		rowData, err := brdr.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return orders, err
		}
		if err == io.EOF && len(rowData) == 0 {
			break
		}
		if order, err := ParseOrder(rowData); err == nil {
			orders = append(orders, order)
		} else {
			log.Println("Parsing order error:", err)
		}
		if err == io.EOF {
			break
		}
	}

	return orders, nil
}

func (p *OrderReader) RemoveCache() error {
	return p.HashStore.Remove()
}
