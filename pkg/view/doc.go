// Package view implements Gated Writer and log writer
package view

import (
	"io"
)

var _ io.Writer = &GatedWriter{}
