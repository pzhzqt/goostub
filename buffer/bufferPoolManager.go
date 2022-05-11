// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package buffer

import (
	"goostub/common"
)

type CallbackType int

const (
	BEFORE CallbackType = iota
	AFTER
)

type bufferpoolCallback func(t CallbackType, pid common.PageID)

type BufferPoolManager struct {
	// TODO
}

func (m *BufferPoolManager) UnpinPage(pid common.PageID, isDirty bool, fn bufferpoolCallback) bool {
	// TODO
}
