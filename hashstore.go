package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
)

type HashPair struct {
	url, chunk []byte // pair of hex hash url and hex hash chunk
}

type HashStore struct {
	fname string // filename
	data  []*HashPair
}

func LoadHashStore(fname string) (*HashStore, error) {
	file, err := os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return &HashStore{fname, nil}, nil
		}
		return nil, err
	}
	defer file.Close()
	dec := json.Decoder(file)
	var data map[string]string
	if err = dec.Decode(&data); err != nil {
		if err == io.EOF {
			return &HashStore{fname, nil}, nil
		}
		return nil, err
	}
	hexdata := make([]*HashPair, 0)
	for url, chunk := range data {
		hex_url, err := hex.DecodeString(url)
		if err != nil {
			// corrupted data
			continue
		}
		hex_chunk, err := hex.DecodeString(chunk)
		if err != nil {
			continue
		}
		hexdata = append(hexdata, &HashPair{hex_url, hex_chunk})
	}
	return &HashStore{fname, hexdata}, nil
}

func (hs *HashStore) Flush() error {
	file, err := os.Create(hs.fname)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make(map[string]string)
	for _, pair := range hs.data {
		data[hex.EncodeToString(pair.url)] =
			hex.EncodeToString(pair.chunk)
	}
	enc := json.NewEncoder(file)
	return enc.Encode(data)
}

func (hs *HashStore) GetHashChunk(rawurl string) ([]byte, bool) {
	if len(rawurl) > 0 {
		hash := md5.New()
		hash.Write([]byte(rawurl))
		url := hash.Sum(nil)
		for _, pair := range hs.data {
			if bytes.Compare(pair.url, url) == 0 {
				return pair.chunk, true
			}
		}
	}
	return nil, false
}

// func (hs *HashStore) ReadHashChunk(rawurl string, r *bufio.Reader) {

// }
