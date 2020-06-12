package types

import (
	"common"
	"github.com/go-kit/kit/log/level"
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

func (t *TinyintType) CastAs(v *Value, id TypeID) (*Value, error) {
    val, ok := v.val.(int8)
    if !ok {
        log.Fatalln("tinyint member function called on non-tinyint value")
    }

    if v.IsNull() {
        return newNullValue(id), nil
    }

    switch (id) {
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
        return NewValue(id, val), nil
    case DECIMAL:
        return NewValue(id, float64(val)), nil
    case VARCHAR:
        return NewValue(id, v.ToString()), nil
    default:
        break
    }

    level.Error(common.Logger).Log()
    return nil, common.NewErrorf(common.INVALID,
        "tinyint is not coercable to %s", TypeIDToString(id))
}
