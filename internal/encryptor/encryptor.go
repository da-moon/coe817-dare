package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"

	"fmt"
	config "github.com/da-moon/coe817-dare/internal/config"
	header "github.com/da-moon/coe817-dare/internal/header"
	log "github.com/da-moon/coe817-dare/pkg/log"
	segment "github.com/da-moon/coe817-dare/internal/segment"
	stacktrace "github.com/palantir/stacktrace"
	"io"
)

// Encryptor ...
type Encryptor struct {
	reader         io.Reader
	writer         io.Writer
	buffer         segment.Segment
	offset         int
	lastByte       byte
	firstRead      bool
	cipher         cipher.AEAD
	randVal        []byte
	sequenceNumber uint32
	finalized      bool
}

// New returns an io.Reader that encrypts everything it reads.
func New(reader io.Reader, writer io.Writer, key []byte) (*Encryptor, error) {
	var err error
	result := &Encryptor{
		reader:    reader,
		writer:    writer,
		buffer:    make(segment.Segment, config.MaxBufferSize),
		firstRead: true,
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
	var randVal [12]byte
	_, err = io.ReadFull(rand.Reader, randVal[:])
	if err != nil {
		return nil, err
	}
	result.randVal = randVal[:]
	return result, nil
}

// Read ...
func (e *Encryptor) Read(p []byte) (int, error) {
	var (
		count int
		err   error
	)
	// cases :
	// 1 - size finalized
	// 2 - not finalized
	if e.firstRead {
		e.firstRead = false
		_, err = io.ReadFull(e.reader, e.buffer[header.HeaderSize:header.HeaderSize+1])
		if err != nil && err != io.EOF {
			return 0, err
		}

		if err == io.EOF {
			e.finalized = true
			return 0, io.EOF
		}
		e.lastByte = e.buffer[header.HeaderSize]
	}
	// log.Debug(fmt.Sprintf("#%d [ENCRYPTOR] finalized=%v firstread=%v lastByte=%v payload_len=%v",
	// 	e.sequenceNumber,
	// 	e.finalized,
	// 	e.firstRead,
	// 	e.lastByte,
	// 	len(p),
	// ),
	// )

	// write the buffered data to p
	if e.offset > 0 {
		remaining := e.buffer.GetLength() - e.offset
		// p(target) doesn't have enough room for
		// remaining data in buffer ... copying all
		// we can , moving offset and retrning ...
		if len(p) < remaining {
			remaining = len(p)
			e.offset += copy(p, e.buffer[e.offset:e.offset+remaining])
			return remaining, nil
		}
		// there is still data left in buffer .
		// copying to p
		count = copy(p, e.buffer[e.offset:e.offset+remaining])
		// updating p
		p = p[remaining:]
		// setting offeset to zero to read more
		e.offset = 0
	}
	if e.finalized {
		return count, io.EOF
	}
	finalize := false
	// as long as reader slice has capacity
	// this for loop would read as long as slice to
	// be populated is larger/equa to encrypted io reader
	// underlying buffer
	for len(p) >= config.MaxBufferSize {
		log.Debug(fmt.Sprintf("#[ENCRYPTOR] seq=%v inside",
			e.sequenceNumber,
		),
		)
		e.buffer[header.HeaderSize] = e.lastByte
		// Reading maximum possible amount
		nn, err := io.ReadFull(
			e.reader,
			e.buffer[header.HeaderSize+1:header.HeaderSize+config.MaxPayloadSize+1],
		)
		if err != nil &&
			err != io.EOF &&
			err != io.ErrUnexpectedEOF {
			err = stacktrace.Propagate(err, "[ERROR] Encryptor failed to read maximum payload from reader")
			return count, err
		}
		// if we are reading less than 64KB , encryptor would seal and finalize
		if err == io.EOF ||
			err == io.ErrUnexpectedEOF {
			finalize = true
			e.seal(p, e.buffer[header.HeaderSize:header.HeaderSize+1+nn], finalize)
			return count + header.HeaderSize + segment.TagSize + 1 + nn, io.EOF
		}
		// saving last read byte for the next burst
		e.lastByte = e.buffer[header.HeaderSize+config.MaxPayloadSize]
		e.seal(p, e.buffer[header.HeaderSize:header.HeaderSize+config.MaxPayloadSize], finalize)
		p = p[config.MaxBufferSize:]
		count += config.MaxBufferSize
	}
	log.Debug(fmt.Sprintf("#[ENCRYPTOR] seq=%v outside len=%v max buf=%v",
		e.sequenceNumber,
		len(p), config.MaxBufferSize,
	),
	)
	if len(p) > 0 {
		e.buffer[header.HeaderSize] = e.lastByte
		nn, err := io.ReadFull(
			e.reader,
			e.buffer[header.HeaderSize+1:header.HeaderSize+config.MaxPayloadSize+1],
		)

		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			err = stacktrace.Propagate(err, "[ERROR] Encryptor failed to read from reader")
			return count, err
		}
		// if we are reading less than 64KB , encryptor would seal and finalize
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			finalize = true
			e.seal(e.buffer, e.buffer[header.HeaderSize:header.HeaderSize+nn+1], finalize)
			if len(p) > e.buffer.GetLength() {
				count += copy(p, e.buffer[:e.buffer.GetLength()])
				return count, io.EOF
			}
		} else {
			// saving last read byte for the next burst
			e.lastByte = e.buffer[header.HeaderSize+config.MaxPayloadSize]
			e.seal(e.buffer, e.buffer[header.HeaderSize:header.HeaderSize+config.MaxPayloadSize], finalize)
		}
		e.offset = copy(p, e.buffer[:len(p)])
		count += e.offset
	}
	return count, nil
}

// Write ...
func (e *Encryptor) Write(p []byte) (int, error) {
	var (
		err      error
		count    int
		finalize = false
	)
	if e.finalized {
		err = stacktrace.NewError("[ERROR] Write to stream after Close is not permitted")
		panic(err)
	}
	if e.offset > 0 {
		remaining := config.MaxPayloadSize - e.offset
		//  buffer == 64 KB
		if len(p) <= remaining {
			e.offset += copy(e.buffer[header.HeaderSize+e.offset:], p)
			return len(p), nil
		}
		count = copy(e.buffer[header.HeaderSize+e.offset:], p[:remaining])
		e.seal(e.buffer, e.buffer[header.HeaderSize:header.HeaderSize+config.MaxPayloadSize], finalize)
		err = flush(e.writer, e.buffer)
		if err != nil {
			return count, err
		}
		p = p[remaining:]
		e.offset = 0
	}
	for len(p) > config.MaxPayloadSize {
		e.seal(e.buffer, p[:config.MaxPayloadSize], finalize)
		err = flush(e.writer, e.buffer)
		if err != nil {
			err = stacktrace.Propagate(err, "[ERROR] Encryptor failed to write to reader")
			return count, err
		}
		p = p[config.MaxPayloadSize:]
		count += config.MaxPayloadSize
	}
	if len(p) > 0 {
		e.offset = copy(e.buffer[header.HeaderSize:], p)
		count += e.offset
	}
	return count, nil
}

// Close ...
func (e *Encryptor) Close() error {
	if e.offset > 0 {
		finalize := true
		e.seal(e.buffer, e.buffer[header.HeaderSize:header.HeaderSize+e.offset], finalize)
		err := flush(e.writer, e.buffer[:header.HeaderSize+e.offset+segment.TagSize])
		if err != nil {
			return err
		}
		e.offset = 0
	}
	closer, ok := e.writer.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}

func (e *Encryptor) seal(dst, src []byte, finalize bool) {
	if e.finalized {
		err := stacktrace.NewError("[ERROR] sealing byte bursts after Close is not permitted")
		panic(err)
	}
	e.finalized = finalize
	h := header.Header(dst[:header.HeaderSize])
	h.SetLength(len(src))
	h.SetNonce(e.randVal, finalize)
	var nonce [header.NonceFieldSize]byte
	copy(nonce[:], h.GetNonce())
	binary.LittleEndian.PutUint32(
		nonce[config.SeqTrackerBit:],
		binary.LittleEndian.Uint32(nonce[config.SeqTrackerBit:])^e.sequenceNumber,
	)

	e.cipher.Seal(dst[header.HeaderSize:header.HeaderSize], nonce[:], src, h.GetAdditionalData())
	e.sequenceNumber++
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
