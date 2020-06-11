package types

import (
	"log"
    "common"
    "github.com/go-kit/kit/log/level"
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
    level.Error(common.Logger).Log("Sqrt not implemented")
    return nil
}

func (t *BaseType) OperateNull(l *Value,r *Value) *Value {
    level.Error(common.Logger).Log("OperateNull not implemented")
    return nil
}

func (t *BaseType) IsZero(v *Value) int8 {
    level.Error(common.Logger).Log("IsZero not implemented")
    return -1
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
    level.Error(common.Logger).Log("Copy not implemented")
    return nil
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

func newInvalidType() Type {
    return &BaseType {
        id: INVALID,
    }
}
