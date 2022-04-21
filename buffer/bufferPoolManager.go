// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package buffer

import (
	"goostub/storage/page"
)

type CallbackType int

const (
	BEFORE CallbackType = iota
	AFTER
)
