package main

import (
	"bytes"
	"crypto/md5"
	"io"
)

// CacheReader reads while does not find chunk which md5 hash sum
// will matches with passed hash
type CacheReader struct {
	r       io.Reader
	delim   []byte // chunk delimiter
	ch_hash []byte // check hash
	chunk   []byte // current chunk
}

// NewCacheReader creates new CacheReader with reader r, delimiter d
// and check hash h
func NewCacheReader(r io.Reader, d, h []byte) *CacheReader {
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

// StackReader is designed to read streams in which chunks of
// information are prepends at the begining such as feeds.
// StackReader saves md5 hash of the first (newest) chunk.
type StackReader struct {
	r     io.Reader
	delim []byte // chunk delimiter
	fline []byte // first (newest) chunk
	flag  bool   // if true first was saved and may be hashed
}

func NewStackReader(r io.Reader, d, h []byte) *StackReader {
	return &StackReader{r, d, nil, false}
}

func (sr *StackReader) Read(p []byte) (n int, err error) {
	n, err = sr.r.Read(p)
	if sr.flag {
		if i := bytes.Index(p, sr.delim); i == -1 {
			sr.fline = append(sr.fline, p...)
		} else {
			sr.fline = append(sr.fline, p[:i]...)
			sr.flag = true
		}
	}
	return
}

func (sr *StackReader) MD5FirstLine() []byte {
	if sr.flag {
		hash := md5.New()
		hash.Write(sr.fline)
		return hash.Sum(nil)
	}
	return nil
}

// HashPoint is pair of md5(url) and md5(chunk) which was readed
// latter
type HashPoint struct {
	url, row string
}

type Cache []*HashPoint

func (c *Cache) Get() string {
	return ""
}
