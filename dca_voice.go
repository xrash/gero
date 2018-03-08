package gero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"time"
)

type DCAVoiceWriter struct {
	io.Writer
	vc                  *discordgo.VoiceConnection
	buffer              chan []byte
	frames              chan []byte
	finishFrameProducer chan bool
	finishFrameSender   chan bool
	running             bool
	paused              bool
	aborted             bool
	onFinish            func()
}

func (vw *DCAVoiceWriter) IsRunning() bool {
	return vw.running
}

func (vw *DCAVoiceWriter) Finish(onFinish func()) {
	vw.onFinish = onFinish
	vw.finishFrameProducer <- true
	vw.finishFrameSender <- true
}

func (vw *DCAVoiceWriter) Pause() {
	vw.paused = true
}

func (vw *DCAVoiceWriter) Resume() {
	vw.paused = false
}

func (vw *DCAVoiceWriter) Abort() {
	vw.aborted = true
}

func NewDCAVoiceWriter(vc *discordgo.VoiceConnection, writesBufferSize, framesBufferSize int) *DCAVoiceWriter {
	return &DCAVoiceWriter{
		vc:                  vc,
		buffer:              make(chan []byte, writesBufferSize),
		frames:              make(chan []byte, framesBufferSize),
		finishFrameProducer: make(chan bool, 1),
		finishFrameSender:   make(chan bool, 1),
	}
}

func (vw *DCAVoiceWriter) Run() error {
	if vw.running {
		return fmt.Errorf("Already running.")
	}

	vw.running = true
	go vw.produceFrames()
	go vw.sendFrames()

	return nil
}

func (vw *DCAVoiceWriter) FramesBuffered() int {
	return len(vw.frames)
}

func (vw *DCAVoiceWriter) Write(p []byte) (int, error) {
	vw.buffer <- p
	return len(p), nil
}

func (vw *DCAVoiceWriter) separateDcaIntoFrames(data []byte) (int, error) {
	var framelen int16
	var err error
	var bytesRead int
	r := bytes.NewBuffer(data)

	for {
		// Get the Opus frame size.
		err = binary.Read(r, binary.LittleEndian, &framelen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			return bytesRead, err
		}

		// Get frame content.
		frame := make([]byte, framelen)
		err = binary.Read(r, binary.LittleEndian, &frame)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			return bytesRead, err
		}

		// int16 + frame
		bytesRead += 2 + int(framelen)

		vw.frames <- frame
	}

	return bytesRead, nil
}

func (vw *DCAVoiceWriter) produceFrames() {
	f := NewDCAFramerizer()
	done := false

outerLoop:
	for {
		select {
		case <-vw.finishFrameProducer:
			done = true
			if len(vw.buffer) == 0 {
				break outerLoop
			}
		case data := <-vw.buffer:
			f.Framerize(data, vw.frames)
			if done && (len(vw.buffer) == 0) {
				break outerLoop
			}
		}
	}
}

func (vw *DCAVoiceWriter) sendFrames() {
	done := false
	defer func() {
		vw.running = false
	}()

outerLoop:
	for {
		if vw.aborted {
			return
		}

		if vw.paused {
			time.Sleep(time.Millisecond * 50)
			continue
		}

		select {
		case <-vw.finishFrameSender:
			done = true
			if (len(vw.buffer) == 0) && (len(vw.frames) == 0) {
				break outerLoop
			}
		case frame := <-vw.frames:
			vw.vc.OpusSend <- frame

			if done && (len(vw.buffer) == 0) && (len(vw.frames) == 0) {
				break outerLoop
			}
		}
	}

	vw.onFinish()
}
