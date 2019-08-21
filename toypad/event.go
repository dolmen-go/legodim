package toypad

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrFrameInvalidLength   = errors.New("invalid frame length")
	ErrFrameInvalidChecksum = errors.New("invalid frame checksum")
)

type inMessage interface {
	fmt.Formatter
}

func parseFrame(frame [32]byte) (payload []byte, err error) {
	l := int(frame[1])
	if l > 31 {
		return nil, ErrFrameInvalidLength
	}
	if frame[2+l] != checkSum(frame[:2+l]) {
		return nil, ErrFrameInvalidChecksum
	}
	return frame[2 : 2+l], nil
}

func parseFrameOld(frame [32]byte) (*Event, error) {
	l := int(frame[1])
	if l > 31 {
		return nil, ErrFrameInvalidLength
	}
	if frame[2+l] != checkSum(frame[:2+l]) {
		return nil, ErrFrameInvalidChecksum
	}
	if frame[0] == 'V' && l > 4 { // 0x56
		return parseEvent(frame[2 : 2+l])
	}
	// TODO responses
	return nil, nil
}

// TagAction represents the kind of event happenning related to a tag on the toypad.
type TagAction uint8

const (
	Add    TagAction = 0 // Add happens when a tag is placed on the toypad.
	Remove TagAction = 1 // Remove happends when a tag is removed from the toypad.
)

func (a TagAction) String() string {
	switch a {
	case Add:
		return "add"
	case Remove:
		return "rmv"
	default:
		return fmt.Sprint(int(a))
	}
}

type Event struct {
	Pad    Pad
	X      uint8 // Unknown data
	Index  uint8
	Action TagAction
	UID    UID
}

func parseEvent(frame []byte) (*Event, error) {
	log.Printf("Event [% X]", frame)
	ev := Event{
		Pad:    Pad(frame[0]),
		X:      frame[1],
		Index:  frame[2],
		Action: TagAction(frame[3]),
	}
	copy(ev.UID[:], frame[4:4+7])
	return &ev, nil
}

func (ev *Event) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if !f.Flag('#') {
			fmt.Fprintf(f, "pad=%d X=%d index=%d action=%s UID=%s", ev.Pad, ev.X, ev.Index, ev.Action, ev.UID)
		} else {
			fmt.Fprintf(f, "%%!%c(%T=)", c, ev)
		}
	default:
		fmt.Fprintf(f, "%%!%c(%T=)", c, ev)
	}
}
