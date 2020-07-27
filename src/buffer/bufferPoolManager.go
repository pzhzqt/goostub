package buffer

import (
	"storage/page"
)

type CallbackType int

const (
	BEFORE CallbackType = iota
	AFTER
)
