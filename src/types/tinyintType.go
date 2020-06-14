package types

import (
    "log"
    "unsafe"
    "strconv"
)

type TinyintType struct {
    IntegerBaseType
}

func newTinyintType() *TinyintType {
    return &TinyintType {
        *newIntegerBaseType(TINYINT),
    }
}

func (t *TinyintType) ToString(v *Value) (string, error) {
    if v.IsNull() {
        return "tinyint_null", nil
    }

    val, ok := v.val.(int8)
    if !ok {
        log.Fatalln("tinyint member function called on non-tinyint value")
    }

    return strconv.FormatInt(int64(val), 10), nil
}

func (t *TinyintType) SerializeTo(v *Value, storage *byte) error {
    val, ok := v.val.(int8)
    if !ok {
        log.Fatalln("tinyint member function called on non-tinyint value")
    }

    *(*int8)(unsafe.Pointer(storage)) = val
    return nil
}

func (t *TinyintType) DeserializeFrom(storage *byte) (*Value, error) {
    val := *(*int8)(unsafe.Pointer(storage))
    return NewValue(t.id, val), nil
}
