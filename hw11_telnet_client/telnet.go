package main

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	conn    net.Conn
	addr    string
	in      io.ReadCloser
	out     io.Writer
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *client) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.addr, t.timeout)
	if err != nil {
		return err
	}

	return nil
}

func (t *client) Close() (err error) {
	if t.conn != nil {
		err = t.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *client) Send() error {
	if t.conn == nil {
		return errors.New("no conn")
	}

	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return err
	}

	return nil
}

func (t *client) Receive() error {
	if t.conn == nil {
		return errors.New("no conn")
	}

	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return err
	}

	return nil
}
