// +build !linux

package toypadlink

import (
	"errors"

	"github.com/karalabe/hid"
)

func List(vendorID, productID uint16) ([]Connect, error) {
	devs := hid.Enumerate(vendorID, productID)
	if len(devs) == 0 {
		return nil, errors.New("no devices found")
	}

	d := make([]Connect, len(devs))
	for i := range devs {
		d[i] = func() (Link, error) {
			return devs[i].Open()
		}
	}
	return d, nil
}
