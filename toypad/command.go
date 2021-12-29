package toypad

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

// Command is a command to send to the toy pad.
type Command struct {
	cmd     cmd
	payload []byte
	cb      func(payload []byte, err error)
}

// Format implements fmt.Formatter.
func (cmd *Command) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if !f.Flag('#') {
			fmt.Fprintf(f, "cmd=%02X payload=[% 02X]", cmd.cmd, cmd.payload)
		}
	default:
		fmt.Fprintf(f, "%%!%c(%T=)", c, cmd)
	}
}

func newCommand(cmd cmd, cb func([]byte, error), args ...interface{}) *Command {
	buf := make([]byte, 0, 32-5)
	b := buf
	var l int
	for _, arg := range args {
		switch arg := arg.(type) {
		case bool:
			if arg {
				b = append(b, 1)
			} else {
				b = append(b, 0)
			}
			b = b[1:]
			l++
		case byte:
			b = append(b, arg)
			b = b[1:]
			l++
		case Pad:
			b = append(b, byte(arg))
			b = b[1:]
			l++
		case uint16:
			b = append(b, 0, 0)
			binary.BigEndian.PutUint16(b, arg)
			b = b[2:]
			l += 2
		case uint32:
			b = append(b, 0, 0, 0, 0)
			binary.BigEndian.PutUint32(b, arg)
			b = b[4:]
			l += 4
		case RGB:
			b = append(b, arg.R, arg.G, arg.B)
			b = b[3:]
			l += 3
		case []byte:
			b = append(b, arg...)
			b = b[len(arg):]
			l += len(arg)
		case string:
			b = append(b, arg...)
			b = b[len(arg):]
			l += len(arg)
		default:
			panic(fmt.Errorf("unsupported type %T", arg))
		}
	}
	return &Command{cmd: cmd, payload: buf[:l], cb: cb}
}

func (c *Command) frame(msgId uint8) []byte {
	buf := [32]byte{}
	buf[0] = 0x55 // Magic host -> portal
	buf[1] = byte(2 + len(c.payload))
	buf[2] = byte(c.cmd)
	buf[3] = msgId
	b := append(buf[0:4], c.payload...)
	b = append(b, checkSum(b[:]))
	return buf[:] // All frames are 32 bytes long
}

func (c *Command) Send(w io.Writer, msgId uint8) error {
	f := c.frame(msgId)
	log.Printf(">> [% X] %[1]q", f[:2+int(f[1])])
	n, err := w.Write(f)
	if err != nil {
		return fmt.Errorf("Write: %w", err)
	}
	if n < len(f) {
		return fmt.Errorf("Write %d/%d: %w", n, len(f), io.ErrShortWrite)
	}
	if n > len(f) {
		// This happens on Windows
		log.Printf("Sent %d bytes, %d written. WTF??", len(f), n)
	}
	return nil
}

/*
func (c *Command) MarshalBinary() ([]byte, error) {
	// TODO
	return nil, nil
}
*/

type cmd byte

const (
	// Source: https://github.com/ags131/node-ld/blob/master/src/lib/constants.js

	cmdInit      cmd = 0xB0
	cmdSeed      cmd = 0xB1
	cmdChallenge cmd = 0xB3

	cmdColor       cmd = 0xC0
	cmdGetColor    cmd = 0xC1
	cmdFade        cmd = 0xC2
	cmdFlash       cmd = 0xC3
	cmdFadeRandom  cmd = 0xC4
	cmdFadeUnknown cmd = 0xC5 // Length: 5
	cmdFadeAll     cmd = 0xC6
	cmdFlashAll    cmd = 0xC7
	cmdColorAll    cmd = 0xC8

	cmdTagList  cmd = 0xD0
	cmdTagRead  cmd = 0xD2
	cmdTagWrite cmd = 0xD3
	cmdTagModel cmd = 0xD4

	cmdConfigPwd    cmd = 0xE1
	cmdConfigEnable cmd = 0xE5
	cmdLEDSQ        cmd = 0xFF
)

func Wake() *Command {
	return newCommand(cmdInit, nil, "(c) LEGO 2014")
}

func SetColor(pad Pad, color RGB) *Command {
	return newCommand(cmdColor, nil, pad, color)
}

func Flash(pad Pad, ticksOn uint8, ticksOff uint8, ticksCount uint8, color RGB) *Command {
	return newCommand(cmdFlash, nil, pad, ticksOn, ticksOff, ticksCount, color)
}

func TagRead(index uint8, pageNum uint8, cb func(status uint8, data []byte, err error)) *Command {
	var cb2 func(payload []byte, err error)
	if cb != nil {
		cb2 = func(payload []byte, err error) {
			if err != nil {
				cb(0, nil, err)
				return
			}
			// payload[0] == 0xF0 or 0xF2 seems to be read failures
			status := payload[0]
			data := payload[1:]
			if payload[0] == 0 {
				log.Printf("\033[1mPages %2d [% X]  %2d [% X]  %2d [% X]  %2d [% X] %q\033[m",
					pageNum, data[:4],
					pageNum+1, data[4:8],
					pageNum+2, data[8:12],
					pageNum+3, data[12:],
					data)
			}
			cb(status, data, nil)
		}
	}
	return newCommand(cmdTagRead, cb2, index, pageNum)
}

func TagModel(index uint8 /* FIXME */, data []byte, cb func(data []byte, err error)) *Command {
	return newCommand(cmdTagModel, cb, data)
}
