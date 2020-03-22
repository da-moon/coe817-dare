package header

import (
	"encoding/binary"
)

const (
	// LengthFieldSize ... number of bits used for storing segment length
	LengthFieldSize = 4
	// NonceFieldSize ... number of bits used for storing nonce
	NonceFieldSize = 12
	// LengthFieldStart ...
	LengthFieldStart = 0
	// NonceFieldStart ...
	NonceFieldStart = LengthFieldStart + LengthFieldSize
	// NonceFieldEnd ...
	NonceFieldEnd = NonceFieldStart + NonceFieldSize
	// HeaderSize ... total header size
	// this is just used for tests
	HeaderSize = LengthFieldSize + NonceFieldSize
)

// Header : encapsulated a byte slice that represents the header
// you can potentially add version by taking
// bytes from size field
type Header []byte

// GetLength ...
func (h Header) GetLength() int {
	slice := h[LengthFieldStart:NonceFieldStart]
	return int(binary.LittleEndian.Uint32(slice)) + 1
}

// SetLength ...
func (h Header) SetLength(length int) {
	// @TODO add tests for correct length
	binary.LittleEndian.PutUint32(h[LengthFieldStart:NonceFieldStart], uint32(length-1))
}

// IsFinal ... checks nonce flag bit (MSB) to see if it was the lastsegment or not
func (h Header) IsFinal() bool {
	return h[NonceFieldStart]|0x7F == 0x7F
}

// GetNonce ...
func (h Header) GetNonce() []byte {
	return h[NonceFieldStart:NonceFieldEnd]
}

// SetNonce ...
func (h Header) SetNonce(nonce []byte, final bool) {
	copy(h[NonceFieldStart:], nonce)
	if !final {
		h.MoreFragments()
		return
	}
	h.FinalFragment()
	return
}

// MoreFragments sets msb of nonce to 1
func (h Header) MoreFragments() {
	//  h[4] or 1000 0000 = 1xxx xxxx
	h[NonceFieldStart] = h[NonceFieldStart] | 0x80
}

// FinalFragment sets msb of nonce to 0
func (h Header) FinalFragment() {
	//  h[4] and 0111 1111 = 0xxx xxxx
	h[NonceFieldStart] = h[NonceFieldStart] & 0x7F
}
