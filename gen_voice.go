package gero

import (
	"fmt"
	"io"
	"os/exec"
)

type GenericVoiceWriter struct {
	io.WriteCloser
	dcaw        *DCAVoiceWriter
	ffmpegStdin io.WriteCloser
	runningFlag bool
}

func NewGenericVoiceWriter(dcaw *DCAVoiceWriter) *GenericVoiceWriter {
	return &GenericVoiceWriter{
		dcaw: dcaw,
	}
}

func (vw *GenericVoiceWriter) Run(onFinish func()) error {
	if vw.runningFlag {
		return fmt.Errorf("Already running.")
	}

	go vw.runCommands(onFinish)

	return nil
}

func (vw *GenericVoiceWriter) Write(p []byte) (int, error) {
	//	fmt.Println("write to gvw", len(p))
	return vw.ffmpegStdin.Write(p)
}

func (vw *GenericVoiceWriter) Close() error {
	return vw.ffmpegStdin.Close()
}

func (vw *GenericVoiceWriter) runCommands(onFinish func()) {
	ffmpeg := exec.Command("ffmpeg", "-i", "-", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	dca := exec.Command("dca")

	ffmpegStdin, err := ffmpeg.StdinPipe()
	if err != nil {
		panic(fmt.Errorf("Error getting stdinpipe from ffmpeg: %v", err))
	}

	ffmpegStdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		panic(fmt.Errorf("Error getting stdoutpipe from ffmpeg: %v", err))
	}

	dca.Stdin = ffmpegStdout
	dca.Stdout = vw.dcaw
	vw.ffmpegStdin = ffmpegStdin

	err = dca.Start()
	if err != nil {
		panic(fmt.Errorf("Error starting dca: %v", err))
	}

	if err := ffmpeg.Start(); err != nil {
		panic(fmt.Errorf("Error starting ffmpeg: %v", err))
	}

	if err := ffmpeg.Wait(); err != nil {
		panic(fmt.Errorf("Error waiting for ffmpeg: %v", err))
	}

	if err := dca.Wait(); err != nil {
		panic(fmt.Errorf("Error waiting dca: %v", err))
	}

	vw.dcaw.Finish(func() {
		vw.runningFlag = false
		onFinish()
	})
}
