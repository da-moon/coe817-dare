package hashsink

import (
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
	md5Hash    hash.Hash
	sha256Hash hash.Hash
}

// NewWriter ...
func NewWriter(
	writer io.Writer,
) *Writer {
	sha256Hash := sha256.New()
	md5Hash := md5.New()
	return &Writer{
		writer:     writer,
		md5Hash:    md5Hash,
		sha256Hash: sha256Hash,
	}
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

	return
}

// MD5 ...
func (w *Writer) MD5() []byte {
	if w.md5Hash != nil {
		return w.md5Hash.Sum(nil)
	}
	return nil
}

// SHA256 ...
func (w *Writer) SHA256() []byte {
	if w.sha256Hash != nil {
		return w.sha256Hash.Sum(nil)
	}
	return nil

}

// MD5HexString ...
func (w *Writer) MD5HexString() string {
	res := w.MD5()
	return hex.EncodeToString(res)
}

// MD5Base64String ...
func (w *Writer) MD5Base64String() string {
	res := w.MD5()
	return base64.StdEncoding.EncodeToString(res)
}

// SHA256HexString ...
func (w *Writer) SHA256HexString() string {
	res := w.SHA256()
	return hex.EncodeToString(res)
}

// SHA256Base64String ...
func (w *Writer) SHA256Base64String() string {
	res := w.SHA256()
	return base64.StdEncoding.EncodeToString(res)
}
