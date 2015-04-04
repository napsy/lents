package lents

import (
	"bytes"
	"errors"
	"io"
	"time"
)

type DataBuffer struct {
	bytes.Buffer
	capacity int
	timeout  time.Duration
}

var (
	ErrFull = errors.New("buffer full") // The data buffer is full
)

const (
	TimeoutDefault time.Duration = 60 // The default timeout value is 60 seconds
)

// Pop reads data from the buffer and puts it in the provided byte slice.
// If the buffer doesn't store the requested size, the function will block
// until sufficient data is available or until the read timeout is reached.
func (buf *DataBuffer) Pop(p []byte) (int, error) {
	var (
		c   chan bool = make(chan bool, 1)
		n   int
		err error
	)
	go func() {
		n, err = buf.Read(p)
		c <- true
	}()
	select {
	case <-time.After(buf.timeout * time.Second):
		return 0, io.EOF
	case <-c:
	}
	return n, err
}

// Push pushes the data into the buffer thus making it available for Pop()
// calls to read them. If the buffer is full, the ErrFull error is returned.
func (buf *DataBuffer) Push(p []byte) (int, error) {
	if len(p)+buf.Len() > buf.capacity {
		return 0, ErrFull
	}
	return buf.Write(p)
}
