package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt" // debug
	"io"
	"os"
)

// CacheReader reads while does not find chunk which md5 hash sum
// will matches with passed checking hash
type CacheReader struct {
	r       io.Reader
	delim   []byte // chunk delimiter
	ch_hash []byte // checking hex hash
	chunk   []byte // current chunk
}

// NewCacheReader creates new CacheReader with reader r, delimiter d
// and check hash h
func NewCacheReader(r io.Reader, d, h []byte) *CacheReader {
	if len(d) == 0 {
		panic("NewCacheReader(): passed nil delimiter")
	}
	return &CacheReader{r, d, h, nil}
}

func (cr *CacheReader) Read(p []byte) (n int, err error) {
	if n, err = cr.r.Read(p); len(cr.delim) > 0 && n > 0 {
		var i, j, k int
		hash := md5.New()
		for {
			k = bytes.Index(p[i:], cr.delim)
			if k == -1 {
				cr.chunk = append(cr.chunk, p[i:]...)
				break
			}
			j += k
			DEB := append(cr.chunk, p[i:j]...) // debug
			fmt.Println(DEB)                   // debug
			hash.Write(DEB)                    // debug
			if bytes.Compare(cr.ch_hash, hash.Sum(nil)) == 0 {
				n = j
				err = io.EOF
				break
			}
			cr.chunk = nil
			j += len(cr.delim)
			i = j
			hash.Reset()
		}
	}
	return
}

// HashPair is object containing pairs: md5 url address and md5 chunk
// which was newest in last parsing. HashPair is "point of last stop"
type HashPair struct {
	url, chunk []byte // pair of md5 url and md5 chunk
}

// HashStore stores urls and chunks on which "feed" was ended at last
// time. HashStore saves info in json file fname
type HashStore struct {
	fname string // file name
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
	dec := json.NewDecoder(file)
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
		//fmt.Printf("get %q\n", rawurl) // debug
		hash := md5.New()
		hash.Write([]byte(rawurl))
		url := hash.Sum(nil)
		//fmt.Printf("get %x\n", url) // debug
		for _, pair := range hs.data {
			if bytes.Compare(pair.url, url) == 0 {
				fmt.Println("found in cache") // debug
				return pair.chunk, true
			}
		}
	}
	return nil, false
}

// SetHashChunk hashes and saves chunk in pair with hashed url
func (hs *HashStore) SetHashChunk(rawurl string, chunk []byte) {
	if len(rawurl) == 0 {
		return
	}
	fmt.Println("check chunk =", chunk) // debug
	// fmt.Printf("set %q\n", rawurl)      // debug
	hash := md5.New()
	hash.Write([]byte(rawurl))
	url := hash.Sum(nil)
	//fmt.Printf("set %x\n", url) // debug
	hash.Reset()
	hash.Write(chunk)
	chunk = hash.Sum(nil)
	for _, pair := range hs.data {
		if bytes.Compare(pair.url, url) == 0 {
			// set existian pair
			pair.chunk = chunk
			return
		}
	}
	// create pair
	hs.data = append(hs.data, &HashPair{url, chunk})
}
