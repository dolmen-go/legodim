package main

import (
	"log"

	"github.com/DizzyDoug/legodim/toypad"
	"github.com/karalabe/hid"
)

func main() {
	devs := hid.Enumerate(toypad.VendorID, toypad.ProductID)
	if len(devs) == 0 {
		log.Fatalln("No devices found.")
	}
	log.Printf("%#v\n", devs)

	dev, err := devs[0].Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer dev.Close()

	tp, err := toypad.NewToyPad(dev, dev)
	if err != nil {
		log.Fatalln(err)
	}

	var count [3]int // count of on each pad
	for {
		select {
		case ev := <-tp.Events:
			log.Printf("Event: %v\n", ev)
			if ev.Action == toypad.Add {
				count[int(ev.Pad-1)]++
				cb := func(uint8, []byte, error) {}
				tp.Send(toypad.TagRead(ev.Index, 0x23, cb))
				tp.Send(toypad.TagRead(ev.Index, 0x27, cb))
				tp.Send(toypad.TagRead(ev.Index, 0x2B, cb))

				// tp.Send(toypad.SetColor(ev.Pad, toypad.RGB{ev.UID[0], ev.UID[1], ev.UID[2]}))
				tp.Send(toypad.Flash(ev.Pad, 5, 10, 50, toypad.RGB{ev.UID[4], ev.UID[5], ev.UID[6]}))
				// tp.Send(toypad.TagModel(ev.Index, nil, nil))
			} else {
				count[int(ev.Pad-1)]--
				if count[int(ev.Pad-1)] == 0 {
					tp.Send(toypad.SetColor(ev.Pad, toypad.RGB{0, 0, 0})) // Switch off
				}
			}
		case err := <-tp.Errors:
			log.Println("Error:", err)
		}
	}

	// Init
	toypad.Wake().Send(dev, 1)

	// Read forever
	for {
		var frame [32]byte
		_, err := dev.Read(frame[:])
		if err != nil {
			log.Println(err)
			break
		}
		l := int(frame[1])
		log.Printf("[% X] %[1]q", frame[:2+l])
	}
}
