package config

import (
	"github.com/da-moon/coe817-dare/internal/header"
	humanize "github.com/dustin/go-humanize"
)

// Consts
const (
	// KeySize ...
	KeySize = 32
	// SeqTrackerBit ...
	SeqTrackerBit = 8

	// MaxPayloadSize : 64KB
	// MaxPayloadSize = 1 << 16
	// MaxPayloadSize = 64 * humanize.KiByte
	MaxPayloadSize = 64 * humanize.Byte
	// TagSize ...
	TagSize = 16
	// Meta
	MetadataSize = header.HeaderSize + TagSize
	// MaxBufferSize ...
	MaxBufferSize = MetadataSize + MaxPayloadSize
	//
	// for test ...
	//
	// MaxDecryptedSize ...
	MaxDecryptedSize = 1 << 48
	// MaxEncryptedSize ...
	MaxEncryptedSize = MaxDecryptedSize + ((MetadataSize) * 1 << 32)
)
