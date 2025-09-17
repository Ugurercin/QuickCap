package recorder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"example.com/internal/common"
	"example.com/internal/config"
)

func findFFmpeg() string {
	exeName := exeName()

	if cwd, err := os.Getwd(); err == nil {
		candidate := filepath.Join(cwd, exeName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidate := filepath.Join(exeDir, exeName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	if p, err := exec.LookPath("ffmpeg"); err == nil {
		return p
	}

	return exeName // last resort, will error if not found
}

func exeName() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

// TODO Arguments should be aware of the operating system -- needs rework
func buildRecordArgs(outputPath string, cfg config.Config) []string {
	out := cfg.Output

	fps := out.FPS
	if fps != 30 && fps != 60 {
		fps = 30 // enforce supported framerate
	}

	switch runtime.GOOS {
	case "darwin":
		switch out.Mode {
		case common.CaptureScreen:
			return []string{
				"-y",
				"-f", "avfoundation",
				"-framerate", strconv.Itoa(fps),
				"-pix_fmt", "uyvy422",
				"-i", "Capture screen 0",
				"-c:v", "libx264",
				"-preset", "ultrafast",
				"-crf", "23",
				"-pix_fmt", "yuv420p",
				outputPath,
			}

		case common.CaptureCamera:
			return []string{
				"-y",
				"-f", "avfoundation",
				"-framerate", strconv.Itoa(fps),
				"-i", "0:none",
				"-c:v", "libx264",
				"-preset", "ultrafast",
				"-crf", "23",
				"-pix_fmt", "yuv420p",
				outputPath,
			}

		case common.CaptureAsOverlay:
			overlay := overlayExpr(out.Position, out.WebcamW, out.WebcamH)
			return []string{
				"-y",
				// Screen input
				"-f", "avfoundation", "-framerate", strconv.Itoa(fps),
				"-pix_fmt", "uyvy422",
				"-i", "Capture screen 0",
				// Camera input
				"-f", "avfoundation", "-framerate", strconv.Itoa(fps),
				"-i", "0:none",
				// Scale & overlay webcam
				"-filter_complex",
				fmt.Sprintf("[1:v] scale=%d:%d [pip]; [0:v][pip] overlay=%s",
					out.WebcamW, out.WebcamH, overlay),
				"-c:v", "libx264",
				"-preset", "ultrafast",
				"-crf", "23",
				"-pix_fmt", "yuv420p",
				outputPath,
			}
		}
	}

	return []string{}

}
func overlayExpr(pos common.OverlayPosition, w, h int) string {
	switch pos {
	case common.PositionTopLeft:
		return "10:10"
	case common.PositionTopRight:
		return fmt.Sprintf("W-%d-10:10", w)
	case common.PositionBottomLeft:
		return fmt.Sprintf("10:H-%d-10", h)
	default: // bottom-right
		return fmt.Sprintf("W-%d-10:H-%d-10", w, h)
	}
}
