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

// Writer ...
type Writer struct {
	writer     io.Writer
	md5sum     []byte
	sha256sum  []byte
	md5Hash    hash.Hash
	sha256Hash hash.Hash
}

// NewWriter ...
func NewWriter(writer io.Writer, md5Hex, sha256Hex string) (*Writer, error) {
	if _, ok := writer.(*Writer); ok {
		return nil, errNestedWriter
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

	return &Writer{
		md5sum:     md5sum,
		sha256sum:  sha256sum,
		writer:     writer,
		md5Hash:    md5Hash,
		sha256Hash: sha256Hash,
	}, nil
}

// Read ...
func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if n > 0 {
		if w.md5Hash != nil {
			w.md5Hash.Write(p[:n])
		}
		if w.sha256Hash != nil {
			w.sha256Hash.Write(p[:n])
		}
	}

	if err == io.EOF {
		if cerr := w.Verify(); cerr != nil {
			return 0, cerr
		}
	}

	return
}

// MD5 ...
func (w *Writer) MD5() []byte {
	return w.md5sum
}

// MD5Current ...
func (w *Writer) MD5Current() []byte {
	if w.md5Hash != nil {
		return w.md5Hash.Sum(nil)
	}
	return nil
}

// SHA256 ...
func (w *Writer) SHA256() []byte {
	return w.sha256sum
}

// MD5HexString ...
func (w *Writer) MD5HexString() string {
	return hex.EncodeToString(w.md5sum)
}

// MD5Base64String ...
func (w *Writer) MD5Base64String() string {
	return base64.StdEncoding.EncodeToString(w.md5sum)
}

// SHA256HexString ...
func (w *Writer) SHA256HexString() string {
	return hex.EncodeToString(w.sha256sum)
}

// Verify ...
func (w *Writer) Verify() error {
	if w.sha256Hash != nil && len(w.sha256sum) > 0 {
		if sum := w.sha256Hash.Sum(nil); !bytes.Equal(w.sha256sum, sum) {
			return SHA256Mismatch{hex.EncodeToString(w.sha256sum), hex.EncodeToString(sum)}
		}
	}
	if w.md5Hash != nil && len(w.md5sum) > 0 {
		if sum := w.md5Hash.Sum(nil); !bytes.Equal(w.md5sum, sum) {
			return BadDigest{hex.EncodeToString(w.md5sum), hex.EncodeToString(sum)}
		}
	}
	return nil
}
