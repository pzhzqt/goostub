// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package types

import (
	"bytes"
	"github.com/go-kit/kit/log/level"
	"goostub/common"
	"math"
)

type Value struct {
	// use interface{} for union type
	val    interface{}
	isNull bool

	// only matters if it's a varchar type
	// since val stores only pointer, this indicates whether
	// this value makes a local copy of the actual data
	manageData bool
	typeID     TypeID
}

/*****************************/
/*******Value Methods********/
/*****************************/

func (v *Value) CheckInteger() bool {
	switch v.GetTypeID() {
	case TINYINT:
		fallthrough
	case SMALLINT:
		fallthrough
	case INTEGER:
		fallthrough
	case BIGINT:
		return true
	}

	return false
}

func (v *Value) CheckComparable(other *Value) bool {
	switch v.GetTypeID() {
	case BOOLEAN:
		id := other.GetTypeID()
		return id == BOOLEAN || id == VARCHAR
	case TINYINT:
		fallthrough
	case SMALLINT:
		fallthrough
	case INTEGER:
		fallthrough
	case BIGINT:
		fallthrough
	case DECIMAL:
		switch other.GetTypeID() {
		case TINYINT:
			fallthrough
		case SMALLINT:
			fallthrough
		case INTEGER:
			fallthrough
		case BIGINT:
			fallthrough
		case DECIMAL:
			fallthrough
		case VARCHAR:
			return true
		}
	case VARCHAR:
		return true
	case TIMESTAMP:
		return other.GetTypeID() == TIMESTAMP
	}

	return false
}

func (v *Value) GetTypeID() TypeID {
	return v.typeID
}

func (v *Value) GetLength() uint32 {
	return uint32(GetInstance(v.typeID).GetLength(v))
}

func (v *Value) GetData() ([]byte, error) {
	return GetInstance(v.typeID).GetData(v)
}

func (v *Value) CastAs(id TypeID) (*Value, error) {
	return GetInstance(v.typeID).CastAs(v, id)
}

// Comparison Method

// a.CompareTo(b) is not symetrical to b.CompareTo(a)
func (v *Value) CompareTo(other *Value) (CmpResult, error) {
	return GetInstance(v.typeID).Compare(v, other)
}

// Other mathematical functions

func (v *Value) Add(other *Value) (*Value, error) {
	return GetInstance(v.typeID).Add(v, other)
}

func (v *Value) Subtract(other *Value) (*Value, error) {
	return GetInstance(v.typeID).Subtract(v, other)
}

func (v *Value) Multiply(other *Value) (*Value, error) {
	return GetInstance(v.typeID).Multiply(v, other)
}

func (v *Value) Divide(other *Value) (*Value, error) {
	return GetInstance(v.typeID).Divide(v, other)
}

func (v *Value) Modulo(other *Value) (*Value, error) {
	return GetInstance(v.typeID).Modulo(v, other)
}

func (v *Value) Min(other *Value) (*Value, error) {
	return GetInstance(v.typeID).Min(v, other)
}

func (v *Value) Max(other *Value) (*Value, error) {
	return GetInstance(v.typeID).Max(v, other)
}

func (v *Value) Sqrt() (*Value, error) {
	return GetInstance(v.typeID).Sqrt(v)
}

func (v *Value) OperateNull(other *Value) (*Value, error) {
	return GetInstance(v.typeID).OperateNull(v, other)
}

func (v *Value) IsZero() (bool, error) {
	return GetInstance(v.typeID).IsZero(v)
}

func (v *Value) IsNull() bool {
	return v.isNull
}

func (v *Value) SerializeTo(storage *bytes.Buffer) error {
	return GetInstance(v.typeID).SerializeTo(v, storage)
}

func (v *Value) DeserializeFrom(storage *bytes.Buffer, id TypeID) (*Value, error) {
	return GetInstance(v.typeID).DeserializeFrom(storage)
}

func (v *Value) ToString() (string, error) {
	return GetInstance(v.typeID).ToString(v)
}

func (v *Value) Copy() *Value {
	return GetInstance(v.typeID).Copy(v)
}

func (v *Value) IsNumeric() bool {
	return v.typeID >= TINYINT && v.typeID <= DECIMAL
}

// stringer for go
func (v Value) String() string {
	s, err := v.ToString()
	if err != nil {
		return err.Error()
	}
	return s
}

/*****************************/
/*****Value Constructors******/
/*****************************/

//clumsy go implementation of constructor overloading
func NewValue(id TypeID, data ...interface{}) *Value {
	if len(data) == 0 {
		return newNullValue(id)
	} else if len(data) == 1 {
		switch data[0].(type) {
		case int8:
			return newValueFromInt8(id, data[0].(int8))
		case int16:
			return newValueFromInt16(id, data[0].(int16))
		case int32:
			return newValueFromInt32(id, data[0].(int32))
		case int64:
			return newValueFromInt64(id, data[0].(int64))
		case int:
			return newValueFromInt64(id, int64(data[0].(int)))
		case uint64:
			return newValueFromUint64(id, data[0].(uint64))
		case float32:
			return newValueFromFloat(id, float64(data[0].(float32)))
		case float64:
			return newValueFromFloat(id, data[0].(float64))
		case []byte:
			return newVarlen(id, data[0].([]byte), true)
		case string:
			return newVarlen(id, ([]byte)(data[0].(string)), true)
		}
	} else if len(data) == 2 {
		slice, sliceOk := data[0].([]byte)
		manageData, boolOk := data[1].(bool)
		if sliceOk && boolOk {
			return newVarlen(id, slice, manageData)
		}
	}

	level.Error(common.Logger).Log("Wrong input format")
	return nil
}

func NewInvalidValue() *Value {
	return newNullValue(INVALID)
}

func CopyValue(other *Value) *Value {
	value := &Value{
		typeID:     other.typeID,
		isNull:     other.isNull,
		manageData: other.manageData,
		val:        other.val,
	}

	if value.typeID == VARCHAR && value.manageData {
		data := other.val.([]byte)
		value.val = make([]byte, len(data))
		copy(value.val.([]byte), data)
	}

	return value
}

func Swap(first *Value, second *Value) {
	//    first.val, second.val = second.val, first.val
	//    first.size, second.size = second.size, first.size
	//    first.manageData, second.manageData = second.manageData, first.manageData
	//    first.typeID, second.typeID = second.typeID, first.typeID
	*first, *second = *second, *first
}

// specific value constructors

func newNullValue(id TypeID) *Value {
	return &Value{
		typeID:     id,
		val:        GetNull(id),
		isNull:     true,
		manageData: false,
	}
}

func newValueFromInt8(id TypeID, i int8) *Value {
	switch id {
	case BOOLEAN:
		fallthrough
	case TINYINT:
		fallthrough
	case SMALLINT:
		fallthrough
	case INTEGER:
		fallthrough
	case BIGINT:
		return newValueFromInt64(id, int64(i))
	}
	level.Error(common.Logger).Log("Invalid Type for 1-byte Value constructor")
	return nil
}

func newValueFromInt16(id TypeID, i int16) *Value {
	switch id {
	case BOOLEAN:
		fallthrough
	case TINYINT:
		fallthrough
	case SMALLINT:
		fallthrough
	case INTEGER:
		fallthrough
	case BIGINT:
		fallthrough
	case TIMESTAMP:
		return newValueFromInt64(id, int64(i))
	}
	level.Error(common.Logger).Log("Invalid Type for 2-byte Value constructor")
	return nil
}

func newValueFromInt32(id TypeID, i int32) *Value {
	switch id {
	case BOOLEAN:
		fallthrough
	case TINYINT:
		fallthrough
	case SMALLINT:
		fallthrough
	case INTEGER:
		fallthrough
	case BIGINT:
		fallthrough
	case TIMESTAMP:
		return newValueFromInt64(id, int64(i))
	}
	level.Error(common.Logger).Log("Invalid Type for 4-byte Value constructor")
	return nil
}

func newValueFromInt64(id TypeID, i int64) *Value {
	value := newNullValue(id)
	var nullval interface{}
	switch id {
	case BOOLEAN:
		value.val = int8(i)
		nullval = GOOSTUB_BOOLEAN_NULL
	case TINYINT:
		value.val = int8(i)
		nullval = GOOSTUB_INT8_NULL
	case SMALLINT:
		value.val = int16(i)
		nullval = GOOSTUB_INT16_NULL
	case INTEGER:
		value.val = int32(i)
		nullval = GOOSTUB_INT32_NULL
	case BIGINT:
		value.val = int64(i)
		nullval = GOOSTUB_INT64_NULL
	case TIMESTAMP:
		value.val = uint64(i)
		nullval = GOOSTUB_TIMESTAMP_NULL
	default:
		level.Error(common.Logger).Log("Invalid Type for 8-byte Value constructor")
		return nil
	}

	if value.val != nullval {
		value.isNull = false
	}

	return value
}

func newValueFromUint64(id TypeID, i uint64) *Value {
	value := newNullValue(id)
	var nullval interface{}
	switch id {
	case BIGINT:
		value.val = int64(i)
		nullval = GOOSTUB_INT64_NULL
	case TIMESTAMP:
		value.val = uint64(i)
		nullval = GOOSTUB_TIMESTAMP_NULL
	default:
		level.Error(common.Logger).Log("Invalid Type for 8-byte Value constructor")
		return nil
	}

	if value.val != nullval {
		value.isNull = false
	}

	return value
}

func newValueFromFloat(id TypeID, d float64) *Value {
	value := newNullValue(id)
	var nullval interface{}
	switch id {
	case DECIMAL:
		value.val = d
		nullval = GOOSTUB_DECIMAL_NULL
	default:
		level.Error(common.Logger).Log("Invalid Type for float Value constructor")
		return nil
	}

	if value.val != nullval {
		value.isNull = false
	}

	return value
}

func newVarlen(id TypeID, data []byte, manageData bool) *Value {
	value := newNullValue(id)
	// switch reserves possibility to add more varlen data type in the future
	switch id {
	case VARCHAR:
		if data == nil {
			value.val = nil
			break
		}

		value.manageData = manageData
		l := len(data)
		value.isNull = false

		if manageData {
			value.val = make([]byte, l)
			copy(value.val.([]byte), data)
		} else {
			value.val = data
		}
	default:
		level.Error(common.Logger).Log("Invalid Type for variable-length Value constructor")
	}

	return value
}

func valMod(x float64, y float64) float64 {
	return x - math.Trunc(x/y)*y
}
