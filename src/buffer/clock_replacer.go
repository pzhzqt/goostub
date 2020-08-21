package buffer

import (
	"goostub/common"
)

type ClockReplacer struct {
	// Student's code
}

// Student: implement everything below

func NewClockReplacer(numPages int64) *ClockReplacer
func (r ClockReplacer) Victim(frameId *common.FrameID) bool
func (r ClockReplacer) Pin(frameId common.FrameID)
func (r ClockReplacer) Unpin(frameId common.FrameID)
func (r ClockReplacer) Size() int64
