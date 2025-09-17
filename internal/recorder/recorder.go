package recorder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"example.com/internal/config"
)

type Recorder struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	file    string
	started time.Time
	active  bool
	cfg     *config.Config
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) Start(outputPath string, cfg config.Config) error {
	args := buildRecordArgs(outputPath, cfg)
	cmd := exec.Command(findFFmpeg(), args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	r.mu.Lock()
	r.cmd = cmd
	r.stdin = stdin
	r.file = outputPath
	r.active = true
	r.mu.Unlock()

	return nil
}

func (r *Recorder) Stop(ctx context.Context) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.active || r.cmd == nil {
		return "", fmt.Errorf("no recording running")
	}

	if r.stdin != nil {
		_, _ = r.stdin.Write([]byte("q\n"))
		_ = r.stdin.Close()
	}

	file := r.file

	done := make(chan error, 1)
	go func(cmd *exec.Cmd) { done <- cmd.Wait() }(r.cmd)

	select {
	case <-ctx.Done():
		if r.cmd.Process != nil {
			_ = r.cmd.Process.Kill()
		}
		r.cleanup()
		return file, fmt.Errorf("timeout waiting for ffmpeg to stop")
	case err := <-done:
		r.cleanup()
		if err != nil && !strings.Contains(err.Error(), "no child processes") {
			return file, fmt.Errorf("ffmpeg exited with error: %w", err)
		}
		return file, nil
	}
}

func (r *Recorder) cleanup() {
	r.active = false
	r.cmd = nil
	r.stdin = nil
	r.file = ""
}

func (r *Recorder) Status() (running bool, file string, elapsed int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cmd == nil {
		return false, "", 0
	}
	return true, r.file, int64(time.Since(r.started).Seconds())
}
