package types

import (
	"common"
	"math"
	"github.com/go-kit/kit/log/level"
)

type IntegerBaseType struct {
    BaseType
}

func newIntegerBaseType(id TypeID) *IntegerBaseType {
    return &IntegerBaseType{
        *newBaseType(id),
    }
}

func (t *IntegerBaseType) Min(l *Value, r *Value) *Value {
    if !t.operandCheck(l, r) {
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
    if !t.operandCheck(l, r) {
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

func (t *IntegerBaseType) Compare(l *Value, r *Value) CmpResult {
    if !t.operandCheck(l, r) {
        return nil
    }

    if l.IsNull() || r.IsNull() {
        return nil
    }

    if !r.CheckInteger() {
        // integer type only handles operations between ints
        return GetInstance(r.typeID).Compare(l, r)
    }

    lval := getValAsBIGINT(l)
    rval := getValAsBIGINT(r)

    var ret int

    if (lval < rval) {
        ret = -1
    } else if lval > rval {
        ret = 1
    } else {
        ret = 0
    }

    return &ret
}

func (t *IntegerBaseType) Add(l *Value, r *Value) *Value {
    if !t.operandCheck(l, r) {
        return nil
    }

    if !r.CheckInteger() {
        // integer type only handles operations between ints
        return GetInstance(r.typeID).Add(l, r)
    }

    if l.IsNull() || r.IsNull() {
        return l.OperateNull(r)
    }

    return t.addInt(l, r)
}

func (t *IntegerBaseType) Subtract(l *Value, r *Value) *Value {
    if !t.operandCheck(l, r) {
        return nil
    }

    if !r.CheckInteger() {
        // integer type only handles operations between ints
        return GetInstance(r.typeID).Subtract(l, r)
    }

    if l.IsNull() || r.IsNull() {
        return l.OperateNull(r)
    }

    return t.subtractInt(l, r)
}

func (t *IntegerBaseType) Multiply(l *Value, r *Value) *Value {
    if !t.operandCheck(l, r) {
        return nil
    }

    if !r.CheckInteger() {
        // integer type only handles operations between ints
        return GetInstance(r.typeID).Multiply(l, r)
    }

    if l.IsNull() || r.IsNull() {
        return l.OperateNull(r)
    }

    return t.multiplyInt(l, r)
}

func (t *IntegerBaseType) Divide(l *Value, r *Value) *Value {
    if !t.operandCheck(l, r) {
        return nil
    }

    if !r.CheckInteger() {
        // integer type only handles operations between ints
        return GetInstance(r.typeID).Divide(l, r)
    }

    if l.IsNull() || r.IsNull() {
        return l.OperateNull(r)
    }

    return t.divideInt(l, r)
}

func (t *IntegerBaseType) Modulo(l *Value, r *Value) *Value {
    if !t.operandCheck(l, r) {
        return nil
    }

    if !r.CheckInteger() {
        // integer type only handles operations between ints
        return GetInstance(r.typeID).Divide(l, r)
    }

    if l.IsNull() || r.IsNull() {
        return l.OperateNull(r)
    }

    return t.moduloInt(l, r)
}

// helper functions

func (t *IntegerBaseType) operandCheck(l *Value, r *Value) bool {
    if !l.CheckInteger() {
        level.Error(common.Logger).Log("Integer member function called from non-integer struct")
        return false
    }

    if !l.CheckComparable(r) {
        level.Error(common.Logger).Log("Not comparable")
        return false
    }

    return true
}

// caller is responsible for null check
func (t *IntegerBaseType) addInt(l *Value, r *Value) *Value {
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
func (t *IntegerBaseType) subtractInt(l *Value, r *Value) *Value {
    nr := NewValue(r.typeID, -getValAsBIGINT(r))
    return t.addInt(l, nr)
}

// caller is responsible for null check
func (t *IntegerBaseType) multiplyInt(l *Value, r *Value) *Value {
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
func (t *IntegerBaseType) divideInt(l *Value, r *Value) *Value {
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
func (t *IntegerBaseType) moduloInt(l *Value, r *Value) *Value {
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
