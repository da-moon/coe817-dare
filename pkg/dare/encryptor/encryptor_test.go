package encryptor_test

import (
	encryptor "github.com/da-moon/dare-cli/pkg/dare/encryptor"
	"io"
)

func init() {
	var _ io.Writer = &encryptor.Writer{}
	var _ io.Reader = &encryptor.Reader{}
}
