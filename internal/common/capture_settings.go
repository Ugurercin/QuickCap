package common

type CaptureMode string

const (
	CaptureScreen    CaptureMode = "screen"
	CaptureCamera    CaptureMode = "camera"
	CaptureAsOverlay CaptureMode = "overlay"
)

type OverlayPosition string

const (
	PositionTopLeft     OverlayPosition = "top-left"
	PositionTopRight    OverlayPosition = "top-right"
	PositionBottomLeft  OverlayPosition = "bottom-left"
	PositionBottomRight OverlayPosition = "bottom-right"
)
