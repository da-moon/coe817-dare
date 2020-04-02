package hashsink

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
)

// Reader ...
type Reader struct {
	reader     io.Reader
	size       int64
	actualSize int64

	md5sum, sha256sum   []byte
	md5Hash, sha256Hash hash.Hash
}

// New ...
func NewReader(reader io.Reader, size int64, md5Hex, sha256Hex string, actualSize int64) (*Reader, error) {
	if _, ok := reader.(*Reader); ok {
		return nil, errNestedReader
	}

	sha256sum, err := hex.DecodeString(sha256Hex)
	if err != nil {
		return nil, SHA256Mismatch{}
	}

	md5sum, err := hex.DecodeString(md5Hex)
	if err != nil {
		return nil, BadDigest{}
	}

	var sha256Hash hash.Hash
	if len(sha256sum) != 0 {
		sha256Hash = sha256.New()
	}
	var md5Hash hash.Hash
	if len(md5sum) != 0 {
		md5Hash = md5.New()
	}
	if size >= 0 {
		reader = io.LimitReader(reader, size)
	}
	return &Reader{
		md5sum:     md5sum,
		sha256sum:  sha256sum,
		reader:     reader,
		size:       size,
		md5Hash:    md5Hash,
		sha256Hash: sha256Hash,
		actualSize: actualSize,
	}, nil
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

	if err == io.EOF {
		if cerr := r.Verify(); cerr != nil {
			return 0, cerr
		}
	}

	return
}

// Size ...
func (r *Reader) Size() int64 { return r.size }

// ActualSize ...
func (r *Reader) ActualSize() int64 { return r.actualSize }

// MD5 ...
func (r *Reader) MD5() []byte {
	return r.md5sum
}

// MD5Current ...
func (r *Reader) MD5Current() []byte {
	if r.md5Hash != nil {
		return r.md5Hash.Sum(nil)
	}
	return nil
}

// SHA256 ...
func (r *Reader) SHA256() []byte {
	return r.sha256sum
}

// MD5HexString ...
func (r *Reader) MD5HexString() string {
	return hex.EncodeToString(r.md5sum)
}

// MD5Base64String ...
func (r *Reader) MD5Base64String() string {
	return base64.StdEncoding.EncodeToString(r.md5sum)
}

// SHA256HexString ...
func (r *Reader) SHA256HexString() string {
	return hex.EncodeToString(r.sha256sum)
}

// Verify ...
func (r *Reader) Verify() error {
	if r.sha256Hash != nil && len(r.sha256sum) > 0 {
		if sum := r.sha256Hash.Sum(nil); !bytes.Equal(r.sha256sum, sum) {
			return SHA256Mismatch{hex.EncodeToString(r.sha256sum), hex.EncodeToString(sum)}
		}
	}
	if r.md5Hash != nil && len(r.md5sum) > 0 {
		if sum := r.md5Hash.Sum(nil); !bytes.Equal(r.md5sum, sum) {
			return BadDigest{hex.EncodeToString(r.md5sum), hex.EncodeToString(sum)}
		}
	}
	return nil
}
