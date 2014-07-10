package main

import (
	"bytes"
	"crypto/md5"
	"io"
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
