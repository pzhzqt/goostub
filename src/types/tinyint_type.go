package types

import (
	"common"
	"github.com/go-kit/kit/log/level"
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

func (t *TinyintType) ToString(v *Value) string {
    if v.IsNull() {
        return "tinyint_null"
    }

    val, ok := v.val.(int8)
    if !ok {
        level.Error(common.Logger).Log("Type Error")
        return ""
    }

    return strconv.FormatInt(int64(val), 10)
}

func (t *TinyintType) SerializeTo(v *Value, storage *byte) {
    val, ok := v.val.(int8)
    if !ok {
        level.Error(common.Logger).Log("Type Error")
        return
    }

    *(*int8)(unsafe.Pointer(storage)) = val
}

func (t *TinyintType) DeserializeFrom(storage *byte) *Value {
    val := *(*int8)(unsafe.Pointer(storage))
    return NewValue(t.id, val)
}

func (t *TinyintType) CastAs(v *Value, id TypeID) *Value {
    val, ok := v.val.(int8)
    if !ok {
        level.Error(common.Logger).Log("Type Error")
        return nil
    }

    if v.IsNull() {
        return newNullValue(id)
    }

    switch (id) {
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
        return NewValue(id, val)
    case DECIMAL:
        return NewValue(id, float64(val))
    case VARCHAR:
        return NewValue(id, v.ToString())
    default:
        break
    }

    level.Error(common.Logger).Log("tinyint is not coercable to this type")
    return nil
}
