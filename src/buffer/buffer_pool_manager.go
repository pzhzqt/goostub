package buffer

import (
    "storage/page"
)

type CallbackType int

const (
    BEFORE CallbackType = iota
    AFTER
)

func boo(x page.Page) {
    x.priv 
}

type bufferPoolCallbackFn func(CallbackType, PageID)

type BufferPoolManager struct {
}
