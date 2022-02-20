package toypadlink

import (
	"errors"

	"github.com/zserge/hid"
)

type hiddev struct {
	dev hid.Device
}

func (d *hiddev) Read(b []byte) (n int, err error) {
	b2, err := d.dev.Read(len(b), 1000)
	if len(b2) != 0 {
		n := copy(b, b2)
		return n, err
	}
	return 0, err
}

func (d *hiddev) Write(b []byte) (n int, err error) {
	n, err = d.dev.Write(b, 1000)
	return
}

func (d *hiddev) Close() error {
	d.dev.Close()
	return nil
}

func List(vendorID, productID uint16) ([]Connect, error) {
	var devs []Connect
	hid.UsbWalk(func(d hid.Device) {
		info := d.Info()
		if info.Vendor != vendorID || info.Product != productID {
			return
		}
		dev := &hiddev{dev: d}
		devs = append(devs, func() (Link, error) {
			return dev, dev.dev.Open()
		})
	})
	if len(devs) == 0 {
		return nil, errors.New("no devices found")
	}

	return devs, nil
}
