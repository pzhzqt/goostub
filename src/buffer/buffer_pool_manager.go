package buffer

import (
    "storage/page"
)

type CallbackType int

const (
    BEFORE CallbackType = iota
    AFTER
)

type bufferPoolCallbackFn func(CallbackType, PageID)

type BufferPoolManager struct {
}
