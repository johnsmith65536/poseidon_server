package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func UnGzip(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	var dst []byte
	var zr *gzip.Reader
	var err error
	if _, err := buf.Write(src); err != nil {
		return nil, err
	}
	if zr, err = gzip.NewReader(&buf); err != nil {
		return nil, err
	}

	if dst, err = ioutil.ReadAll(zr); err != nil {
		return nil, err
	}

	if err = zr.Close(); err != nil {
		return nil, err
	}
	return dst, nil
}
