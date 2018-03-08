package gero

import (
	"bytes"
	"encoding/binary"
	"io"
)

type DCAFramerizer struct {
	leftover []byte
}

func NewDCAFramerizer() *DCAFramerizer {
	return &DCAFramerizer{
		leftover: make([]byte, 0),
	}
}

func (f *DCAFramerizer) Framerize(in []byte, out chan []byte) error {
	merged := append(f.leftover, in...)
	bytesRead, err := f.separateIntoFrames(merged, out)
	if err != nil {
		return err
	}
	f.leftover = merged[bytesRead:]
	return nil
}

func (f *DCAFramerizer) separateIntoFrames(in []byte, out chan []byte) (int, error) {
	var framelen int16
	var err error
	var bytesRead int
	b := bytes.NewBuffer(in)

	for {
		// Get the Opus frame size.
		err = binary.Read(b, binary.LittleEndian, &framelen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			return bytesRead, err
		}

		// Get frame content.
		frame := make([]byte, framelen)
		err = binary.Read(b, binary.LittleEndian, &frame)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			return bytesRead, err
		}

		// int16 + frame
		bytesRead += 2 + int(framelen)

		out <- frame
	}

	return bytesRead, nil
}
