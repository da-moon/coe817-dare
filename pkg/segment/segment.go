package segment

import (
	header "github.com/da-moon/coe817-dare/pkg/header"
)

const (
	// TagSize ...
	TagSize = 16
)

// Segment ...
type Segment []byte

// Header ...
func (s Segment) Header() header.Header {
	headerEndBit := header.HeaderSize
	return header.Header(s[:headerEndBit])
}

// Data ...
func (s Segment) Data() []byte {
	startBit := header.HeaderSize
	rawDataLength := s.Header().GetLength()
	endBit := startBit + rawDataLength
	return s[startBit:endBit]
}

// GetCiphertext ...
func (s Segment) GetCiphertext() []byte {
	return s[header.HeaderSize:s.GetLength()]
}

// GetLength ...
func (s Segment) GetLength() int {
	rawDataLength := s.Header().GetLength()
	return header.HeaderSize + TagSize + rawDataLength
}
