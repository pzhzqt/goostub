// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package buffer

import (
	"goostub/common"
)

type Replacer interface {
	// choose a victim and remove it, assign the frame id to *frameID
	// return true if victim is found, false otherwise
	Victim(frameId *common.FrameID) bool

	// pin a frame
	// whatever is pinned won't be victimized until unpinned
	Pin(frameId common.FrameID)

	// unpin a frame, indicating it can now be victimized
	Unpin(frameId common.FrameID)

	// return the number of elements that can be victimized
	Size() int64
}
