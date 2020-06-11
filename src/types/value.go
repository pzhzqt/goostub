package types

import (
    "log"
    "common"
    "github.com/go-kit/kit/log/level"
    "math"
)

type Value struct {
    // use interface{} for union type
    val         interface{}
    size        uint32

    // only matters if it's a varchar type
    // since val stores only pointer, this indicates whether
    // this value makes a local copy of the actual data
    manageData  bool
    typeID      TypeID
}

/*****************************/
/*******Value Methods********/
/*****************************/

func (v *Value) CheckInteger() bool {
    switch (v.GetTypeID()) {
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
        return true
    default:
        break
    }

    return false
}

func (v *Value) CheckComparable(other *Value) bool {
    switch (v.GetTypeID()) {
    case BOOLEAN:
        id := other.GetTypeID()
        return id == BOOLEAN || id == VARCHAR
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
    case DECIMAL:
        switch (other.GetTypeID()) {
        case TINYINT:
        case SMALLINT:
        case INTEGER:
        case BIGINT:
        case DECIMAL:
        case VARCHAR:
            return true
        default:
            break
        }
        break
    case VARCHAR:
        return true
    default:
        break
    }

    return false
}

func (v *Value) GetTypeID() TypeID {
    return v.typeID
}

func (v *Value) GetLength() uint32 {
    return GetInstance(v.typeID).GetLength(v)
}

func (v *Value) GetData() []byte {
    return GetInstance(v.typeID).GetData(v)
}

func (v *Value) CastAs(id TypeID) *Value {
    return GetInstance(v.typeID).CastAs(v, id)
}

// Comparison Method

func (v *Value) CompareTo(other *Value) CmpResult {
    return GetInstance(v.typeID).Compare(v, other)
}

// Other mathematical functions

func (v *Value) Add(other *Value) *Value {
    return GetInstance(v.typeID).Add(v, other)
}

func (v *Value) Subtract(other *Value) *Value {
    return GetInstance(v.typeID).Subtract(v, other)
}

func (v *Value) Multiply(other *Value) *Value {
    return GetInstance(v.typeID).Multiply(v, other)
}

func (v *Value) Divide(other *Value) *Value {
    return GetInstance(v.typeID).Divide(v, other)
}

func (v *Value) Modulo(other *Value) *Value {
    return GetInstance(v.typeID).Modulo(v, other)
}

func (v *Value) Min(other *Value) *Value {
    return GetInstance(v.typeID).Min(v, other)
}

func (v *Value) Max(other *Value) *Value {
    return GetInstance(v.typeID).Max(v, other)
}

func (v *Value) Sqrt() *Value {
    return GetInstance(v.typeID).Sqrt(v)
}

func (v *Value) OperateNull(other *Value) *Value {
    return GetInstance(v.typeID).OperateNull(v, other)
}

func (v *Value) IsZero() int8 {
    return GetInstance(v.typeID).IsZero(v)
}

func (v *Value) IsNull() bool {
    return v.size == GOOSTUB_VALUE_NULL
}

func (v *Value) SerializeTo(storage *byte) {
    GetInstance(v.typeID).SerializeTo(v, storage)
}

func (v *Value) DeserializeFrom(storage *byte, id TypeID) *Value {
    return GetInstance(v.typeID).DeserializeFrom(storage)
}

func (v *Value) ToString() string {
    return GetInstance(v.typeID).ToString(v)
}

func (v *Value) Copy() *Value {
    return GetInstance(v.typeID).Copy(v)
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
        case int:
            return newValueFromInt32(id, data[0].(int32))
        case int64:
            return newValueFromInt64(id, data[0].(int64))
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
        default:
            break
        }
    } else if len(data) == 2 {
        slice, sliceOk := data[0].([]byte)
        manageData, boolOk := data[1].(bool)
        if sliceOk && boolOk {
            return newVarlen(id, slice, manageData)
        }
    }

    log.Fatalln("Wrong input format")
    return nil
}

func NewInvalidValue() *Value {
    return newNullValue(INVALID)
}

func CopyValue(other *Value) *Value {
    value := &Value {
        typeID     : other.typeID,
        size       : other.size,
        manageData : other.manageData,
        val        : other.val,
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
    return &Value {
        typeID :     id,
        size :       GOOSTUB_VALUE_NULL,
        manageData : false,
    }
}

func newValueFromInt8(id TypeID, i int8) *Value {
    switch (id) {
    case BOOLEAN:
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
        return newValueFromInt64(id, int64(i))
    default:
        break
    }
    level.Error(common.Logger).Log("Invalid Type for 1-byte Value constructor")
    return nil
}

func newValueFromInt16(id TypeID, i int16) *Value {
    switch (id) {
    case BOOLEAN:
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
    case TIMESTAMP:
        return newValueFromInt64(id, int64(i))
    default:
        break
    }
    level.Error(common.Logger).Log("Invalid Type for 2-byte Value constructor")
    return nil
}

func newValueFromInt32(id TypeID, i int32) *Value {
    switch (id) {
    case BOOLEAN:
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
    case TIMESTAMP:
        return newValueFromInt64(id, int64(i))
    default:
        break
    }
    level.Error(common.Logger).Log("Invalid Type for 4-byte Value constructor")
    return nil
}

func newValueFromInt64(id TypeID, i int64) *Value {
    value := newNullValue(id)
    var nullval interface{}
    switch (id) {
    case BOOLEAN:
        value.val = int8(i)
        nullval = GOOSTUB_BOOLEAN_NULL
        break
    case TINYINT:
        value.val = int8(i)
        nullval = GOOSTUB_INT8_NULL
        break
    case SMALLINT:
        value.val = int16(i)
        nullval = GOOSTUB_INT16_NULL
        break
    case INTEGER:
        value.val = int32(i)
        nullval = GOOSTUB_INT32_NULL
        break
    case BIGINT:
        value.val = int64(i)
        nullval = GOOSTUB_INT64_NULL
        break
    case TIMESTAMP:
        value.val = uint64(i)
        nullval = GOOSTUB_TIMESTAMP_NULL
        break
    default:
        level.Error(common.Logger).Log("Invalid Type for 8-byte Value constructor")
        return nil
    }

    if value.val != nullval {
        value.size = 0
    }

    return value
}

func newValueFromUint64(id TypeID, i uint64) *Value {
    value := newNullValue(id)
    var nullval interface{}
    switch (id) {
    case BIGINT:
        value.val = int64(i)
        nullval = GOOSTUB_INT64_NULL
        break
    case TIMESTAMP:
        value.val = uint64(i)
        nullval = GOOSTUB_TIMESTAMP_NULL
        break
    default:
        level.Error(common.Logger).Log("Invalid Type for 8-byte Value constructor")
        return nil
    }

    if value.val != nullval {
        value.size = 0
    }

    return value
}

func newValueFromFloat(id TypeID, d float64) *Value {
    value := newNullValue(id)
    var nullval interface{}
    switch (id) {
    case DECIMAL:
        value.val = d
        nullval = GOOSTUB_DECIMAL_NULL
        break
    default:
        level.Error(common.Logger).Log("Invalid Type for float Value constructor")
        return nil
    }

    if value.val != nullval {
        value.size = 0
    }

    return value
}

func newVarlen(id TypeID, data []byte, manageData bool) *Value {
    value := newNullValue(id)
    switch (id) {
    case VARCHAR:
        if data == nil {
            value.val = nil
            break
        }

        value.manageData = manageData
        l := len(data)
        value.size = uint32(l)

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
    return x - math.Trunc(x/y) * y
}

