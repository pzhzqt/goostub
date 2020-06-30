package types

import (
	"common"
	"github.com/go-kit/kit/log/level"
	"log"
	"math"
)

type IntegerBaseType struct {
	BaseType
}

func newIntegerBaseType(id TypeID) *IntegerBaseType {
	return &IntegerBaseType{
		*newBaseType(id),
	}
}

// Unfortunately there's no macro in go, otherwise this file could be a lot more DRY

func (t *IntegerBaseType) Min(l, r *Value) (*Value, error) {
	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	res, err := t.Compare(l, r)
	if err != nil {
		return nil, err
	}
	if res < 0 {
		return l.Copy(), nil
	}

	return r.Copy(), nil
}

func (t *IntegerBaseType) Max(l, r *Value) (*Value, error) {
	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	res, err := t.Compare(l, r)
	if err != nil {
		return nil, err
	}
	if res > 0 {
		return l.Copy(), nil
	}

	return r.Copy(), nil
}

func (t *IntegerBaseType) Compare(l *Value, r *Value) (CmpResult, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return 0, err
	}

	if l.IsNull() || r.IsNull() {
		return 0, common.NewError(common.INVALID,
			"Null Value is not comparable")
	}

	if !r.CheckInteger() {
		// integer type only handles operations between ints
		return GetInstance(r.typeID).Compare(l, r)
	}

	lval := getValAsBIGINT(l)
	rval := getValAsBIGINT(r)

	var ret CmpResult

	if lval < rval {
		ret = -1
	} else if lval > rval {
		ret = 1
	} else {
		ret = 0
	}

	return ret, nil
}

func (t *IntegerBaseType) Add(l *Value, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	if !r.CheckInteger() {
		// integer type only handles operations between ints
		return GetInstance(r.typeID).Add(l, r)
	}

	id := l.typeID
	if id < r.typeID {
		id = r.typeID
	}

	x := getValAsBIGINT(l)
	y := getValAsBIGINT(r)

	if (x > 0 && y > math.MaxInt64-x) || (x < 0 && y < math.MinInt64-x) {
		// Overflow 64 bits
		return nil, common.NewError(common.OUT_OF_RANGE, "Result out of range")
	}

	res := x + y

	if intOverflow(res, id) {
		return nil, common.NewError(common.OUT_OF_RANGE, "Result out of range")
	}

	return NewValue(id, res), nil
}

// we could use Add to implement substract, however we don't do this to get rid of function call overhead
func (t *IntegerBaseType) Subtract(l *Value, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	if !r.CheckInteger() {
		// integer type only handles operations between ints
		return GetInstance(r.typeID).Add(l, r)
	}

	id := l.typeID
	if id < r.typeID {
		id = r.typeID
	}

	x := getValAsBIGINT(l)
	y := getValAsBIGINT(r)

	if (x > 0 && y < x-math.MaxInt64) || (x < 0 && y > x-math.MinInt64) {
		// Overflow 64 bits
		return nil, common.NewError(common.OUT_OF_RANGE, "Result out of range")
	}

	res := x - y

	if intOverflow(res, id) {
		return nil, common.NewError(common.OUT_OF_RANGE, "Result out of range")
	}

	return NewValue(id, res), nil
}

func (t *IntegerBaseType) Multiply(l *Value, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	if !r.CheckInteger() {
		// integer type only handles operations between ints
		return GetInstance(r.typeID).Add(l, r)
	}

	id := l.typeID
	if id < r.typeID {
		id = r.typeID
	}

	x := getValAsBIGINT(l)
	y := getValAsBIGINT(r)

	if (x > 0 && (y > math.MaxInt64/x || y < math.MinInt64/x)) || (x < 0 && (y < math.MaxInt64/x || y > math.MinInt64/x)) {
		// Overflow 64 bits
		return nil, common.NewError(common.OUT_OF_RANGE, "Result out of range")
	}

	res := x * y

	if intOverflow(res, id) {
		return nil, common.NewError(common.OUT_OF_RANGE, "Result out of range")
	}

	return NewValue(id, res), nil
}

func (t *IntegerBaseType) Divide(l *Value, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	if !r.CheckInteger() {
		// integer type only handles operations between ints
		return GetInstance(r.typeID).Add(l, r)
	}

	id := l.typeID
	if id < r.typeID {
		id = r.typeID
	}

	x := getValAsBIGINT(l)
	y := getValAsBIGINT(r)

	if y == 0 {
		return nil, common.NewError(common.DIVIDE_BY_ZERO,
			"Divide by zero")
	}

	return NewValue(id, x/y), nil
}

func (t *IntegerBaseType) Modulo(l *Value, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	if !r.CheckInteger() {
		// integer type only handles operations between ints
		return GetInstance(r.typeID).Add(l, r)
	}

	id := l.typeID
	if id < r.typeID {
		id = r.typeID
	}

	x := getValAsBIGINT(l)
	y := getValAsBIGINT(r)

	if y == 0 {
		return nil, common.NewError(common.DIVIDE_BY_ZERO,
			"Divide by zero")
	}

	return NewValue(id, x%y), nil
}

func (t *IntegerBaseType) CastAs(v *Value, id TypeID) (*Value, error) {
	if v.IsNull() {
		return newNullValue(id), nil
	}

	val := getValAsBIGINT(v)

	switch id {
	case TINYINT:
	case SMALLINT:
	case INTEGER:
	case BIGINT:
		return NewValue(id, val), nil
	case DECIMAL:
		return NewValue(id, float64(val)), nil
	case VARCHAR:
		s, err := v.ToString()
		if err != nil {
			return nil, err
		}
		return NewValue(id, s), nil
	default:
		break
	}

	return nil, common.NewErrorf(common.INVALID,
		"%s is not coercable to %s", TypeIDToString(v.typeID), TypeIDToString(id))
}

// helper functions

func (t *IntegerBaseType) operandCheck(l *Value, r *Value) error {
	if !l.CheckInteger() {
		log.Fatalln("Integer member function called on non-integer value")
	}

	if !l.CheckComparable(r) {
		return common.NewErrorf(common.MISMATCH_TYPE,
			"Mismatched types of %s and %s", TypeIDToString(l.typeID), TypeIDToString(r.typeID))
	}

	return nil
}

func getValAsBIGINT(v *Value) int64 {
	if !v.CheckInteger() {
		log.Fatalln("getValAsBIGINT should only be called on integer value")
	}

	switch v.GetTypeID() {
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

	log.Fatalln("this shouldn't print out")
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
