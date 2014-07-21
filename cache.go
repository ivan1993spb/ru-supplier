package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
)

const _HASH_STORE_FILE_NAME = "cache.json"

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

// Read reads and checks chunks
func (cr *CacheReader) Read(p []byte) (n int, err error) {
	if n, err = cr.r.Read(p); len(cr.delim) > 0 && n > 0 {
		var i, j, k int
		hash := md5.New()
		for {
			k = bytes.Index(p[i:n], cr.delim)
			if k == -1 {
				cr.chunk = append(cr.chunk, p[i:n]...)
				break
			}
			j += k
			hash.Write(append(cr.chunk, p[i:j]...))
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

func LoadHashStoreSimple() *HashStore {
	hs, err := LoadHashStore(_HASH_STORE_FILE_NAME)
	if hs == nil {
		if err != nil {
			log.Fatal("cannot load hashstore:", err)
		} else {
			panic("hashstore is nil")
		}
	}
	if err != nil {
		log.Println("hashstore:", err)
	}
	return hs
}

func LoadHashStore(fname string) (hs *HashStore, err error) {
	if len(fname) == 0 {
		panic("hashstore: invalid file name")
	}
	var file *os.File
	hs = &HashStore{fname, nil}
	file, err = os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	var data map[string]string
	if err = dec.Decode(&data); err != nil {
		if err == io.EOF {
			err = nil
		}
		return
	}
	hs.data = make([]*HashPair, 0)
	for url, chunk := range data {
		hex_url, err := hex.DecodeString(url)
		if err != nil {
			continue
		}
		hex_chunk, err := hex.DecodeString(chunk)
		if err != nil {
			continue
		}
		hs.data = append(hs.data, &HashPair{hex_url, hex_chunk})
	}
	return
}

func (hs *HashStore) Save() error {
	if len(hs.data) == 0 {
		hs.Remove()
		return nil
	}
	file, err := os.Create(hs.fname)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make(map[string]string)
	for _, pair := range hs.data {
		data[hex.EncodeToString(pair.url)] = hex.EncodeToString(pair.chunk)
	}
	return json.NewEncoder(file).Encode(data)
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

// SetHashChunk hashes and saves chunk in pair with hashed url
func (hs *HashStore) SetHashChunk(rawurl string, chunk []byte) {
	if len(rawurl) == 0 {
		return
	}
	hash := md5.New()
	hash.Write([]byte(rawurl))
	url := hash.Sum(nil)
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
	// create pair if new url
	hs.data = append(hs.data, &HashPair{url, chunk})
}

// Remove removes all cache
func (hs *HashStore) Remove() error {
	hs.data = nil
	return os.Remove(hs.fname)
}
