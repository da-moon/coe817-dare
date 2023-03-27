package decryptor_test

import (
	decryptor "github.com/da-moon/dare-cli/pkg/dare/decryptor"
	"io"
)

func init() {
	var _ io.Writer = &decryptor.Writer{}
	var _ io.Reader = &decryptor.Reader{}
}
