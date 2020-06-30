package types

import (
	"common"
	"log"
	"strconv"
	"unsafe"
)

type DecimalType struct {
	BaseType
}

func newDecimalType() *DecimalType {
	return &DecimalType{
		*newBaseType(DECIMAL),
	}
}

func (t *DecimalType) Min(l, r *Value) (*Value, error) {
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

func (t *DecimalType) Max(l, r *Value) (*Value, error) {
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

func (t *DecimalType) Compare(l, r *Value) (CmpResult, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return 0, err
	}

	if l.IsNull() || r.IsNull() {
		return 0, common.NewError(common.INVALID,
			"Null Value is not comparable")
	}

	lval := getValAsDecimal(l)
	rval := getValAsDecimal(r)

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

func (t *DecimalType) Add(l, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	// prolly should check overflow but since double overflow is highly unlikely, we might as well not add overhead
	return NewValue(DECIMAL, getValAsDecimal(l)+getValAsDecimal(r)), nil
}

func (t *DecimalType) Subtract(l, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	// prolly should check overflow but since double overflow is highly unlikely, we might as well not add overhead
	return NewValue(DECIMAL, getValAsDecimal(l)-getValAsDecimal(r)), nil
}

func (t *DecimalType) Multiply(l, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	// prolly should check overflow but since double overflow is highly unlikely, we might as well not add overhead
	return NewValue(DECIMAL, getValAsDecimal(l)*getValAsDecimal(r)), nil
}

func (t *DecimalType) Divide(l, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	if z, err := r.IsZero(); z || err != nil {
		return nil, common.NewError(common.DIVIDE_BY_ZERO,
			"Divide by zero")
	}

	// prolly should check overflow but since double overflow is highly unlikely, we might as well not add overhead
	return NewValue(DECIMAL, getValAsDecimal(l)/getValAsDecimal(r)), nil
}

func (t *DecimalType) Modulo(l, r *Value) (*Value, error) {
	err := t.operandCheck(l, r)
	if err != nil {
		return nil, err
	}

	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	if z, err := r.IsZero(); z || err != nil {
		return nil, common.NewError(common.DIVIDE_BY_ZERO,
			"Divide by zero")
	}

	// prolly should check overflow but since double overflow is highly unlikely, we might as well not add overhead
	return NewValue(DECIMAL, valMod(getValAsDecimal(l), getValAsDecimal(r))), nil
}

func (t *DecimalType) CastAs(v *Value, id TypeID) (*Value, error) {
	if v.IsNull() {
		return newNullValue(id), nil
	}

	val := getValAsDecimal(v)
	var msg string = "Numeric Value out of range"

	switch id {
	case TINYINT:
		if val > float64(GOOSTUB_INT8_MAX) || val < float64(GOOSTUB_INT8_MIN) {
			return nil, common.NewError(common.OUT_OF_RANGE, msg)
		}
		return NewValue(id, int8(val)), nil
	case SMALLINT:
		if val > float64(GOOSTUB_INT16_MAX) || val < float64(GOOSTUB_INT16_MIN) {
			return nil, common.NewError(common.OUT_OF_RANGE, msg)
		}
		return NewValue(id, int16(val)), nil
	case INTEGER:
		if val > float64(GOOSTUB_INT32_MAX) || val < float64(GOOSTUB_INT32_MIN) {
			return nil, common.NewError(common.OUT_OF_RANGE, msg)
		}
		return NewValue(id, int32(val)), nil
	case BIGINT:
		if val > float64(GOOSTUB_INT64_MAX) || val < float64(GOOSTUB_INT64_MIN) {
			return nil, common.NewError(common.OUT_OF_RANGE, msg)
		}
		return NewValue(id, int64(val)), nil
	case DECIMAL:
		return v.Copy(), nil
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

func (t *DecimalType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "decimal_null", nil
	}

	val, ok := v.val.(float64)
	if !ok {
		log.Fatalln("decimal member function called on non-decimal value")
	}

	return strconv.FormatFloat(val, 'G', -1, 64), nil
}

func (t *DecimalType) SerializeTo(v *Value, storage *byte) error {
	val, ok := v.val.(float64)
	if !ok {
		log.Fatalln("decimal member function called on non-decimal value")
	}

	*(*float64)(unsafe.Pointer(storage)) = val
	return nil
}

func (t *DecimalType) DeserializeFrom(storage *byte) (*Value, error) {
	val := *(*float64)(unsafe.Pointer(storage))
	return NewValue(t.id, val), nil
}

// helper functions

func (t *DecimalType) operandCheck(l, r *Value) error {
	if l.GetTypeID() != DECIMAL && r.GetTypeID() != DECIMAL {
		log.Fatalln("Decimal function called on non-decimal values")
	}

	if !l.CheckComparable(r) {
		return common.NewErrorf(common.MISMATCH_TYPE,
			"Mismatched types of %s and %s", TypeIDToString(l.typeID), TypeIDToString(r.typeID))
	}

	return nil
}

func getValAsDecimal(v *Value) float64 {
	if !v.IsNumeric() {
		log.Fatalln("getValAsDecimal should only be called on integer value")
	}

	switch v.GetTypeID() {
	case TINYINT:
		return float64(v.val.(int8))
	case SMALLINT:
		return float64(v.val.(int16))
	case INTEGER:
		return float64(v.val.(int32))
	case BIGINT:
		return float64(v.val.(int64))
	case DECIMAL:
		return v.val.(float64)
	default:
		break
	}

	log.Fatalln("this shouldn't print out")
	return GOOSTUB_DECIMAL_NULL
}
