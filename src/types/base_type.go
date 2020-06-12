package types

import (
    "common"
    "github.com/go-kit/kit/log/level"
	"log"
    "math"
)

type BaseType struct {
    id TypeID
}

func (t *BaseType) IsCoercableFrom(id TypeID)bool {
    switch (t.id) {
    case INVALID:
        return false
    case BOOLEAN:
        return true
    case TINYINT:
    case SMALLINT:
    case INTEGER:
    case BIGINT:
    case DECIMAL:
        switch (id) {
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
        switch (id) {
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
func (t *BaseType) Compare(l *Value, r *Value) CmpResult {
    level.Error(common.Logger).Log("CompareEquals not implemented")
    return nil
}

// Math Functions
func (t *BaseType) Add(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("Add not implemented")
    return nil
}

func (t *BaseType) Subtract(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("Subtract not implemented")
    return nil
}

func (t *BaseType) Multiply(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("Multiply not implemented")
    return nil
}

func (t *BaseType) Divide(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("Divide not implemented")
    return nil
}

func (t *BaseType) Modulo(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("Modulo not implemented")
    return nil
}

func (t *BaseType) Min(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("Min not implemented")
    return nil
}

func (t *BaseType) Max(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("Max not implemented")
    return nil
}

func (t *BaseType) Sqrt(v *Value) *Value {
    var val float64

    switch (v.typeID) {
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
        level.Error(common.Logger).Log("Can't take square root of non-numeric type")
        return nil
    }

    if v.IsNull() {
        return NewValue(DECIMAL, GOOSTUB_DECIMAL_NULL)
    }

    if val < 0 {
        level.Error(common.Logger).Log("Can't take square root of a negative number")
        return nil
    }

    return NewValue(DECIMAL, math.Sqrt(val))
}

func (t *BaseType) OperateNull(l *Value, r *Value) *Value {
    if !l.IsNumeric() || !r.IsNumeric() {
        level.Error(common.Logger).Log("OperateNull not implemented for types")
        return nil
    }

    if l.typeID < r.typeID {
        return newNullValue(r.typeID)
    } else {
        return newNullValue(l.typeID)
    }
}

func (t *BaseType) IsZero(v *Value) int8 {
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
        level.Error(common.Logger).Log("IsZero not implemented for this type")
        return -1
    }

    if val == 0 {
        return 1
    }

    return 0
}

func (t *BaseType) IsInlined(v *Value) int8 {
    level.Error(common.Logger).Log("IsInlined not implemented")
    return -1
}

func (t *BaseType) ToString(v *Value) string {
    level.Error(common.Logger).Log("ToString not implemented")
    return ""
}

func (t *BaseType) SerializeTo(v *Value, storage *byte) {
    level.Error(common.Logger).Log("SerializeTo not implemented")
}

func (t *BaseType) DeserializeFrom(storage *byte) *Value {
    level.Error(common.Logger).Log("DeserializeFrom not implemented")
    return nil
}

func (t *BaseType) Copy(v *Value) *Value {
    return CopyValue(v)
}

func (t *BaseType) CastAs(v *Value, id TypeID) *Value {
    level.Error(common.Logger).Log("CastAs not implemented")
    return nil
}

// raw variable-length data
func (t *BaseType) GetData(v *Value) []byte {
    level.Error(common.Logger).Log("GetData not implemented")
    return nil
}

// length of variable-length data
func (t *BaseType) GetLength(v *Value) uint32 {
    level.Error(common.Logger).Log("GetLenth not implemented")
    return 0
}

func newBaseType(tid TypeID) *BaseType {
    return &BaseType {
        id: tid,
    }
}

func newInvalidType() *BaseType {
    return newBaseType(INVALID)
}
