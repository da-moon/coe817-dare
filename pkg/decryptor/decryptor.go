package decryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"io"

	config "github.com/da-moon/coe817-dare/pkg/config"
	errors "github.com/da-moon/coe817-dare/pkg/errors"
	header "github.com/da-moon/coe817-dare/pkg/header"
	segment "github.com/da-moon/coe817-dare/pkg/segment"
	stacktrace "github.com/palantir/stacktrace"
)

// Decryptor ...
type Decryptor struct {
	reader         io.Reader
	writer         io.Writer
	buffer         segment.Segment
	finalized      bool
	sequenceNumber uint32
	cipher         cipher.AEAD
	offset         int
}

// New returns an io.Reader decrypts everything it reads.
func New(reader io.Reader, writer io.Writer, key []byte) (*Decryptor, error) {
	result := &Decryptor{

		reader: reader,
		writer: writer,
		buffer: make(segment.Segment, config.MaxBufferSize),
	}
	if len(key) != config.KeySize {
		err := stacktrace.NewError("[ERROR] Encryptor cannot be initialized due to invalid key size")
		return nil, err
	}
	aes256, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	result.cipher, err = cipher.NewGCM(aes256)
	if err != nil {
		return nil, err
	}
	return result, nil

}

// Read ...
func (d *Decryptor) Read(p []byte) (n int, err error) {
	if d.offset > 0 {
		remaining := len(d.buffer.Data()) - d.offset
		if len(p) < remaining {
			n = copy(p, d.buffer.Data()[d.offset:d.offset+len(p)])
			d.offset += n
			return n, nil
		}
		n = copy(p, d.buffer.Data()[d.offset:])
		p = p[remaining:]
		d.offset = 0
	}
	for len(p) >= config.MaxPayloadSize {
		nn, err := io.ReadFull(d.reader, d.buffer)
		if err == io.EOF && !d.finalized {
			err = errors.ErrUnexpectedEOF
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Read failed because it reached EOF without getting final data burst")
			return n, err
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			// err = stacktrace.Propagate(err, "[ERROR] Decryptor Read failed because it reached EOF or reading from reader failed")
			return n, err
		}
		err = d.metadata(p, d.buffer[:nn])
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Read failed because it could not initialize metadata for the sequence")

			return n, err
		}
		p = p[len(d.buffer.Data()):]
		n += len(d.buffer.Data())
	}
	if len(p) > 0 {
		nn, err := io.ReadFull(d.reader, d.buffer)
		if err == io.EOF && !d.finalized {
			err = errors.ErrUnexpectedEOF
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Read failed because it reached EOF without getting final data burst")
			return n, err
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			// err = stacktrace.Propagate(err, "[ERROR] Decryptor Read failed because it reached EOF or reading from reader failed")
			return n, err
		}
		err = d.metadata(d.buffer[header.HeaderSize:], d.buffer[:nn])
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Read failed because it could not initialize metadata for the sequence")
			return n, err
		}
		payload := d.buffer.Data()
		if len(p) < len(payload) {
			d.offset = copy(p, payload[:len(p)])
			n += d.offset
		} else {
			n += copy(p, payload)
		}
	}
	return n, nil
}

// WriteAt ...
func (d *Decryptor) WriteAt(p []byte, off int64) (n int, err error) {
	d.offset = int(off)
	return d.Write(p)
}

// Write ...
func (d *Decryptor) Write(p []byte) (n int, err error) {
	if d.offset > 0 {
		remaining := header.HeaderSize + config.MaxPayloadSize + segment.TagSize - d.offset
		if len(p) < remaining {
			d.offset += copy(d.buffer[d.offset:], p)
			return len(p), nil
		}
		n = copy(d.buffer[d.offset:], p[:remaining])
		plaintext := d.buffer[header.HeaderSize : header.HeaderSize+config.MaxPayloadSize]
		err = d.metadata(plaintext, d.buffer)
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Write failed because it could not initialize metadata for the sequence")
			return n, err
		}
		err = flush(d.writer, plaintext)
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Write failed because it could not flush data to underlying io.Writer")
			return n, err
		}
		p = p[remaining:]
		d.offset = 0
	}
	for len(p) >= config.MaxBufferSize {
		plaintext := d.buffer[header.HeaderSize : header.HeaderSize+config.MaxPayloadSize]
		err = d.metadata(plaintext, p[:config.MaxBufferSize])
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Write failed because it could not initialize metadata for the sequence")
			return n, err
		}

		err = flush(d.writer, plaintext)
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Decryptor Write failed because it could not flush data to underlying io.Writer")
			return n, err
		}
		p = p[config.MaxBufferSize:]
		n += config.MaxBufferSize
	}
	if len(p) > 0 {
		if d.finalized {
			return n, errors.ErrUnexpectedData
		}
		d.offset = copy(d.buffer[:], p)
		n += d.offset
	}
	return n, nil
}

// Close ...
func (d *Decryptor) Close() error {
	if d.offset > 0 {
		// the payload must always be greater than 0
		if d.offset <= header.HeaderSize+segment.TagSize {
			err := errors.ErrInvalidPayloadSize
			err = stacktrace.Propagate(err, "[ERROR] Could not close decryptor because current offset (%v) is larger than sum of header size constant (%v) and Tag size constant (%v)", d.offset, header.HeaderSize, segment.TagSize)
			return err
		}
		destination := d.buffer[header.HeaderSize : d.offset-segment.TagSize]
		source := d.buffer[:d.offset]
		err := d.metadata(destination, source)
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Could not close decryptor because it could not initialize metadata for the sequence")
			return err
		}
		err = flush(
			d.writer,
			d.buffer[header.HeaderSize:d.offset-segment.TagSize],
		)
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Could not close decryptor because it could not flush data to underlying io.Writer")
			return err
		}
		d.offset = 0
	}
	closer, ok := d.writer.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}

// metadata ...
func (d *Decryptor) metadata(dst, src []byte) error {
	if d.finalized {

		return errors.ErrUnexpectedData
	}
	if len(src) <= header.HeaderSize+segment.TagSize {
		err := errors.ErrInvalidPayloadSize
		err = stacktrace.Propagate(err, "[ERROR] Could not generate metadata for decryptor because current source length (%v) is lower than or equal to sum of header size constant (%v) and Tag size constant (%v)", len(src), header.HeaderSize, segment.TagSize)
		return err
	}
	h := segment.Segment(src).Header()
	if len(src) != header.HeaderSize+segment.TagSize+h.GetLength() {
		err := errors.ErrInvalidPayloadSize
		err = stacktrace.Propagate(err, "[ERROR] Could not generate metadata for decryptor because current source length (%v) is not equal to the sum of header size constant (%v) and Tag size constant (%v) and header size (%v)", len(src), header.HeaderSize, segment.TagSize, h.GetLength())
		return err
	}
	if !h.IsFinal() && h.GetLength() != config.MaxPayloadSize {
		err := errors.ErrInvalidPayloadSize
		err = stacktrace.Propagate(err, "[ERROR] Could not generate metadata for decryptor because unfinalized header length (%v) is not equal to the sum of max payload size constant (%v)", h.GetLength(), config.MaxPayloadSize)
		return err
	}
	refNonce := h.GetNonce()
	// refNonce := d.header.GetNonce()
	if h.IsFinal() {
		d.finalized = true
		// refNonce[0] |= config.HeaderFinalFlag
		refNonce[0] = refNonce[0] & 0x7F

	}
	if subtle.ConstantTimeCompare(h.GetNonce(), refNonce[:]) != 1 {
		return errors.ErrNonceMismatch
	}
	var nonce [header.NonceFieldSize]byte
	copy(nonce[:], h.GetNonce())
	binary.LittleEndian.PutUint32(
		nonce[config.SeqTrackerBit:],
		binary.LittleEndian.Uint32(nonce[config.SeqTrackerBit:])^d.sequenceNumber,
	)
	cipher := d.cipher
	// ciphertext := src[header.HeaderSize : header.HeaderSize+header.GetLength()+segment.TagSize]
	// ciphertext := src[header.HeaderSize:segment.Segment(src).GetLength()]
	ciphertext := segment.Segment(src).GetCiphertext()
	_, err := cipher.Open(
		dst[:0],
		nonce[:],
		ciphertext,
		h.GetAdditionalData(),
	)
	if err != nil {
		err = stacktrace.Propagate(err, errors.ErrAuthentication.Error())
		err = stacktrace.Propagate(err, "[ERROR] Could not generate metadata for decryptor becuase cipher failed with decrypting and authenticating ciphertext")
		return err
	}
	d.sequenceNumber++
	return nil
}

func flush(w io.Writer, p []byte) error {
	n, err := w.Write(p)
	if err != nil {
		return err
	}
	if n != len(p) {
		return io.ErrShortWrite
	}
	return nil
}
