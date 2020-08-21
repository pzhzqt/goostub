package buffer

import (
	"goostub/storage/page"
)

type CallbackType int

const (
	BEFORE CallbackType = iota
	AFTER
)
