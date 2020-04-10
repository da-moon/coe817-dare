package hashsink

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
)

// Reader ...
type Reader struct {
	reader io.Reader

	md5Hash    hash.Hash
	sha256Hash hash.Hash
}

// New ...
func NewReader(
	reader io.Reader,
	size int64,

) *Reader {

	sha256Hash := sha256.New()
	md5Hash := md5.New()
	if size >= 0 {
		reader = io.LimitReader(reader, size)
	}
	return &Reader{
		reader:     reader,
		md5Hash:    md5Hash,
		sha256Hash: sha256Hash,
	}
}

// Read ...
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if n > 0 {
		if r.md5Hash != nil {
			r.md5Hash.Write(p[:n])
		}
		if r.sha256Hash != nil {
			r.sha256Hash.Write(p[:n])
		}
	}

	return
}

// MD5 ...
func (r *Reader) MD5() []byte {
	if r.md5Hash != nil {
		return r.md5Hash.Sum(nil)
	}
	return nil
}

// SHA256 ...
func (r *Reader) SHA256() []byte {
	if r.sha256Hash != nil {
		return r.sha256Hash.Sum(nil)
	}
	return nil

}

// MD5HexString ...
func (r *Reader) MD5HexString() string {
	res := r.MD5()
	return hex.EncodeToString(res)
}

// MD5Base64String ...
func (r *Reader) MD5Base64String() string {
	res := r.MD5()
	return base64.StdEncoding.EncodeToString(res)
}

// SHA256HexString ...
func (r *Reader) SHA256HexString() string {
	res := r.SHA256()
	return hex.EncodeToString(res)
}

// SHA256Base64String ...
func (r *Reader) SHA256Base64String() string {
	res := r.SHA256()
	return base64.StdEncoding.EncodeToString(res)
}
