package toypad

import (
	"encoding/hex"
	"io"
	"log"
	"sync"
)

const (
	VendorID  = 0x0e6f
	ProductID = 0x0241
)

// Pad is one area (or all) of the Toy Pad.
type Pad uint8

const (
	AllPads   Pad = 0
	CenterPad Pad = 1
	LeftPad   Pad = 2
	RightPad  Pad = 3
)

// RGB is a color.
type RGB struct {
	R, G, B uint8
}

type UID [7]byte

func (uid UID) String() string {
	return hex.EncodeToString(uid[:])
}

type tagSlot struct {
	pad Pad
	uid UID
}

type ToyPad struct {
	mu sync.Mutex
	w  io.Writer

	msgId uint8
	cb    [256]func([]byte, error)
	// sendQueue []*Command

	tags [12]tagSlot

	debug bool

	events chan *Event
	errors chan error
	Events <-chan *Event
	Errors <-chan error
}

func NewToyPad(r io.Reader, w io.Writer) (*ToyPad, error) {
	events := make(chan *Event, 2)
	errors := make(chan error)
	toypad := ToyPad{
		w:      w,
		msgId:  0,
		events: events,
		Events: events,
		errors: errors,
		Errors: errors,
	}
	// Init
	if err := toypad.Send(Wake()); err != nil {
		return nil, err
	}

	go toypad.readLoop(r)
	return &toypad, nil
}

func (tp *ToyPad) Send(cmd *Command) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.msgId++
	tp.cb[tp.msgId] = cmd.cb
	return cmd.Send(tp.w, tp.msgId)
}

func (tp *ToyPad) readLoop(r io.Reader) {
	for {
		var frame [32]byte
		_, err := r.Read(frame[:])
		if err != nil {
			log.Println(err)
			break
		}
		// log.Printf("[% X]", frame[:])
		payload, err := parseFrame(frame)
		if err != nil {
			tp.errors <- err
			log.Printf("[% X] %v", frame[:], err)
			continue
		}
		l := int(frame[1])
		log.Printf("<< [% X] %[1]q", frame[:2+l])
		switch frame[0] {
		case 'V': // Event
			ev, err := parseEvent(payload)
			if err != nil {
				tp.errors <- err
			} else {
				tp.mu.Lock()
				if ev.Action == Add {
					tp.tags[int(ev.Index)] = tagSlot{pad: ev.Pad, uid: ev.UID}
				} else {
					tp.tags[int(ev.Index)] = tagSlot{} // Clear the slot
				}
				tp.mu.Unlock()
				tp.events <- ev
			}
		case 'U': // Reply to a request
			log.Printf("Reply to msg %d [% X]", payload[0], payload[1:])
			msgId := int(payload[0])
			cb := tp.cb[msgId]
			tp.cb[msgId] = nil
			if cb != nil {
				cb(payload[1:], nil)
			}
		default:
			log.Printf("Unknown frame type %02X", frame[0])
		}
	}
}
