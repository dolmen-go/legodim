package toypadlink

import "io"

type Link interface {
	io.Reader
	io.Writer
	io.Closer
}

type Connect func() (Link, error)
