package util

import (
	"github.com/jedib0t/go-pretty/progress"
	"time"
)

type ProgressBar struct {
	pw progress.Writer
}

func NewProgressBar() ProgressBar {
	pw := progress.NewWriter()
	pw.SetTrackerLength(10)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.SetStyle(progress.StyleBlocks)
	pw.SetMessageWidth(50)
	pw.SetNumTrackersExpected(7)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 50)
	pw.Style().Colors = progress.StyleColorsExample

	return ProgressBar{pw: pw}
}

func (pb *ProgressBar) Render() {
	go pb.pw.Render()
}

func (pb *ProgressBar) Start(msg string) *progress.Tracker {
	t := progress.Tracker{Message: msg, Total: 1}
	pb.pw.AppendTracker(&t)

	return &t
}

func (pb *ProgressBar) StopRendering() {
	for pb.pw.IsRenderInProgress() {
		if pb.pw.LengthActive() == 0 {
			pb.pw.Stop()
		}

		time.Sleep(time.Millisecond * 100)
	}
}
