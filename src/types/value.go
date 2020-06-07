package types

import (
    "unsafe"
)

type Value struct {
    val unsafe.Pointer
}

func NewValue(id TypeID, data ...interface{}) *Value {
    if len(data) == 0 {
        return newNullValue(id)
    }
}

// specific value constructors
func newNullValue(id TypeID) *Value {
}
