package config

import (
	"github.com/da-moon/coe817-dare/pkg/header"
)

// Consts
const (
	// KeySize ...
	KeySize = 32
	// MaxPayloadSize ...
	MaxPayloadSize = 1 << 16
	// TagSize ...
	TagSize = 16
	// Meta
	MetadataSize = header.HeaderSize + TagSize
	// MaxPackageSize ...
	MaxPackageSize = MetadataSize + MaxPayloadSize
	// MaxDecryptedSize ...
	MaxDecryptedSize = 1 << 48
	// MaxEncryptedSize ...
	MaxEncryptedSize = MaxDecryptedSize + ((MetadataSize) * 1 << 32)

	// MaxBufferSize ...
	MaxBufferSize = MetadataSize + MaxPayloadSize
)
