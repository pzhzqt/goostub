package types

import (
	"common"
	"math"
	"github.com/go-kit/kit/log/level"
)

type IntegerBaseType struct {
    BaseType
}

func newIntegerBaseType(id TypeID) Type {
    return newBaseType(id)
}

func (t *IntegerBaseType) Min(l *Value, r *Value) *Value {
    if !l.CheckInteger() {
        level.Error(common.Logger).Log("Integer member function called from non-integer struct")
        return nil
    }

    if !l.CheckComparable(r) {
        level.Error(common.Logger).Log("Not comparable")
        return nil
    }

    if l.IsNull() || r.IsNull() {
        return l.OperateNull(r)
    }

    if *l.CompareTo(r) < 0 {
        return l.Copy()
    }

    return r.Copy()
}

func (t *IntegerBaseType) Max(l *Value, r *Value) *Value {
    if !l.CheckInteger() {
        level.Error(common.Logger).Log("Integer member function called from non-integer struct")
        return nil
    }

    if !l.CheckComparable(r) {
        level.Error(common.Logger).Log("Not comparable")
        return nil
    }

    if l.IsNull() || r.IsNull() {
        return l.OperateNull(r)
    }

    if *l.CompareTo(r) > 0 {
        return l.Copy()
    }

    return r.Copy()
}

// caller is responsible for null check
func (t *IntegerBaseType) addValue(l *Value, r *Value) *Value {
    id := l.typeID
    if id < r.typeID {
        id = r.typeID
    }

    x := getValAsBIGINT(l)
    y := getValAsBIGINT(r)

    if (x > 0 && y > math.MaxInt64 - x) || (x < 0 && y < math.MinInt64 - x) {
    // Overflow 64 bits
        level.Debug(common.Logger).Log("Overflow")
        return nil
    }

    res := x + y

    if intOverflow(res, id) {
        return nil
    }

    return NewValue(id, res)
}

// caller is responsible for null check
func (t *IntegerBaseType) subtractValue(l *Value, r *Value) *Value {
    nr := NewValue(r.typeID, -getValAsBIGINT(r))
    return t.addValue(l, nr)
}

// caller is responsible for null check
func (t *IntegerBaseType) multiplyValue(l *Value, r *Value) *Value {
    id := l.typeID
    if id < r.typeID {
        id = r.typeID
    }

    x := getValAsBIGINT(l)
    y := getValAsBIGINT(r)

    if (x > 0 && (y > math.MaxInt64/x || y < math.MinInt64/x)) || (x < 0 && (y < math.MaxInt64/x || y > math.MinInt64/x)) {
    // Overflow 64 bits
        level.Debug(common.Logger).Log("Overflow")
        return nil
    }

    res := x * y

    if intOverflow(res, id) {
        return nil
    }

    return NewValue(id, res)
}

// caller is responsible for null check
func (t *IntegerBaseType) divideValue(l *Value, r *Value) *Value {
    id := l.typeID
    if id < r.typeID {
        id = r.typeID
    }

    x := getValAsBIGINT(l)
    y := getValAsBIGINT(r)

    if y == 0 {
        level.Error(common.Logger).Log("Divide by zero")
        return nil
    }

    return NewValue(id, x/y)
}

// caller is responsible for null check
func (t *IntegerBaseType) moduloValue(l *Value, r *Value) *Value {
    id := l.typeID
    if id < r.typeID {
        id = r.typeID
    }

    x := getValAsBIGINT(l)
    y := getValAsBIGINT(r)

    if y == 0 {
        level.Error(common.Logger).Log("Divide by zero")
        return nil
    }

    return NewValue(id, x % y)
}

// helper functions

func getValAsBIGINT(v *Value) int64 {
    switch (v.GetTypeID()) {
    case TINYINT:
        return int64(v.val.(int8))
    case SMALLINT:
        return int64(v.val.(int16))
    case INTEGER:
        return int64(v.val.(int32))
    case BIGINT:
        return v.val.(int64)
    default:
        break
    }
    return GOOSTUB_INT64_NULL
}

func intOverflow(v int64, id TypeID) bool {
    if v > getValAsBIGINT(GetMaxValue(id)) {
        level.Debug(common.Logger).Log("Overflow")
        return true
    }

    if v < getValAsBIGINT(GetMinValue(id)) {
        level.Debug(common.Logger).Log("Overflow")
        return true
    }

    return false
}
