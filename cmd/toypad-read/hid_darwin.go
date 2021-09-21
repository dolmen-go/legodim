package main

import (
	"errors"
	"io"

	"github.com/karalabe/hid"
)

type hidconn interface {
	io.Reader
	io.Writer
	io.Closer
}

type hidopen func() (hidconn, error)

func connect(vendorID, productID uint16) ([]hidopen, error) {
	devs := hid.Enumerate(vendorID, productID)
	if len(devs) == 0 {
		return nil, errors.New("no devices found")
	}

	d := make([]hidopen, len(devs))
	for i := range devs {
		d[i] = func() (hidconn, error) {
			return devs[i].Open()
		}
	}
	return d, nil
}
