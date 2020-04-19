package utils

import (
	"bytes"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func EucKrToUtf8(s string) string {
	var buffers bytes.Buffer
	tr := transform.NewWriter(&buffers, korean.EUCKR.NewDecoder())
	defer func() { _ = tr.Close() }()
	_, _ = tr.Write([]byte(s))
	return buffers.String()
}
