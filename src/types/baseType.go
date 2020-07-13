package types

import (
	"bytes"
	"common"
	"log"
	"math"
)

type BaseType struct {
	id TypeID
}

func (t *BaseType) IsCoercableFrom(id TypeID) bool {
	switch t.id {
	case INVALID:
		return false
	case BOOLEAN:
		return true
	case TINYINT:
	case SMALLINT:
	case INTEGER:
	case BIGINT:
	case DECIMAL:
		switch id {
		case TINYINT:
		case SMALLINT:
		case INTEGER:
		case BIGINT:
		case DECIMAL:
			return true
		default:
			return false
		}
	case TIMESTAMP:
		return id == VARCHAR || id == TIMESTAMP
	case VARCHAR:
		switch id {
		case INVALID:
			return false
		default:
			return true
		}
	default:
		return id == t.id
	}

	// stupidity check
	log.Fatalln("If this prints out, there's a code error")
	return false
}

func (t *BaseType) GetTypeID() TypeID {
	return t.id
}

// Comparisons
func (t *BaseType) Compare(l *Value, r *Value) (CmpResult, error) {
	return 0, common.NewError(common.NOT_IMPLEMENTED,
		"Compare not implemented")
}

// Math Functions
func (t *BaseType) Add(l *Value, r *Value) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"Add not implemented")
}

func (t *BaseType) Subtract(l *Value, r *Value) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"Subtract not implemented")
}

func (t *BaseType) Multiply(l *Value, r *Value) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"Multiply not implemented")
}

func (t *BaseType) Divide(l *Value, r *Value) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"Divide not implemented")
}

func (t *BaseType) Modulo(l *Value, r *Value) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"Modulo not implemented")
}

func (t *BaseType) Min(l *Value, r *Value) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"Min not implemented")
}

func (t *BaseType) Max(l *Value, r *Value) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"Max not implemented")
}

func (t *BaseType) Sqrt(v *Value) (*Value, error) {
	var val float64

	switch v.typeID {
	case TINYINT:
		val = float64(v.val.(int8))
		break
	case SMALLINT:
		val = float64(v.val.(int16))
		break
	case INTEGER:
		val = float64(v.val.(int32))
	case BIGINT:
		val = float64(v.val.(int64))
	case DECIMAL:
		val = v.val.(float64)
	default:
		return nil, common.NewErrorf(common.INCOMPATIBLE_TYPE,
			"Can't take square root of type %s", TypeIDToString(v.typeID))
	}

	if v.IsNull() {
		return NewValue(DECIMAL, GOOSTUB_DECIMAL_NULL), nil
	}

	if val < 0 {
		return nil, common.NewError(common.DECIMAL,
			"Can't take square root of a negative number")
	}

	return NewValue(DECIMAL, math.Sqrt(val)), nil
}

func (t *BaseType) OperateNull(l *Value, r *Value) (*Value, error) {
	if !l.IsNumeric() || !r.IsNumeric() {
		return nil, common.NewError(common.INCOMPATIBLE_TYPE,
			"operate_null only works for numeric types")
	}

	if l.typeID < r.typeID {
		return newNullValue(r.typeID), nil
	} else {
		return newNullValue(l.typeID), nil
	}
}

func (t *BaseType) IsZero(v *Value) (bool, error) {
	var val float64
	switch v.val.(type) {
	case int8:
		val = float64(v.val.(int8))
		break
	case int16:
		val = float64(v.val.(int16))
		break
	case int32:
		val = float64(v.val.(int32))
		break
	case int64:
		val = float64(v.val.(int64))
		break
	case float64:
		val = v.val.(float64)
		break
	default:
		return false, common.NewError(common.INCOMPATIBLE_TYPE,
			"This type has no zero value")
	}

	if val == 0 {
		return true, nil
	}

	return false, nil
}

func (t *BaseType) IsInlined(v *Value) (bool, error) {
	return false, common.NewError(common.NOT_IMPLEMENTED,
		"IsInlined not implemented")
}

func (t *BaseType) ToString(v *Value) (string, error) {
	return "", common.NewError(common.NOT_IMPLEMENTED,
		"ToString not implemented")
}

func (t *BaseType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	return common.NewError(common.NOT_IMPLEMENTED,
		"SerializeTo not implemented")
}

func (t *BaseType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"DeserializeFrom not implemented")
}

func (t *BaseType) Copy(v *Value) *Value {
	return CopyValue(v)
}

func (t *BaseType) CastAs(v *Value, id TypeID) (*Value, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"CastAs not implemented")
}

// raw variable-length data
func (t *BaseType) GetData(v *Value) ([]byte, error) {
	return nil, common.NewError(common.NOT_IMPLEMENTED,
		"GetData not implemented")
}

// length of variable-length data
func (t *BaseType) GetLength(v *Value) (uint32, error) {
	return 0, common.NewError(common.NOT_IMPLEMENTED,
		"GetLength not implemented")
}

func newBaseType(tid TypeID) *BaseType {
	return &BaseType{
		id: tid,
	}
}

func newInvalidType() *BaseType {
	return newBaseType(INVALID)
}
