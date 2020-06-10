package types

import (
	"errors"
	"log"
    "common"
    "github.com/go-kit/kit/log/level"
)

type TypeID int

const (
    INVALID TypeID = iota
    BOOLEAN
    TINYINT
    SMALLINT
    INTEGER
    BIGINT
    DECIMAL
    VARCHAR
    TIMESTAMP
)

type CmpBool int

const (
    CmpFalse CmpBool = iota
    CmpTrue
    CmpNull
)

type Type interface {
    IsCoercableFrom(TypeID)bool
    GetTypeID()TypeID

    // Comparisons
    CompareEquals(*Value, *Value)CmpBool
    CompareNotEquals(*Value, *Value)CmpBool
    CompareLessThan(*Value, *Value)CmpBool
    CompareLessThanEquals(*Value, *Value)CmpBool
    CompareGreaterThan(*Value, *Value)CmpBool
    CompareGreaterThanEquals(*Value, *Value)CmpBool

    // Math Functions
    Add(*Value, *Value)*Value
    Subtract(*Value, *Value)*Value
    Multiply(*Value, *Value)*Value
    Divide(*Value, *Value)*Value
    Modulo(*Value, *Value)*Value
    Min(*Value, *Value)*Value
    Max(*Value, *Value)*Value
    Sqrt(*Value)*Value
    OperateNull(*Value, *Value)*Value
    IsZero(*Value)bool
    // Is the data in the struct storage or has indirection
    IsInlined(*Value)bool
    ToString(*Value)string
    SerializeTo(*Value, []byte)
    DeserializeFrom([]byte)*Value
    Copy(*Value)*Value
    CastAs(*Value,TypeID)*Value
    // raw variable length data
    GetData(*Value)[]byte
    GetLength(*Value)uint32
}

type GenericType struct {
    id TypeID
}

func (t *GenericType) IsCoercableFrom(id TypeID)bool {
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

func newInvalidType()Type
func newBooleanType()Type
func newTinyintType()Type
func newSmallintType()Type
func newIntegerType()Type
func newBigintType()Type
func newDecimalType()Type
func newVarcharType()Type
func newTimestampType()Type

var typeConstructors = [...]func()Type {
    newInvalidType,
    newBooleanType,
    newTinyintType,
    newSmallintType,
    newIntegerType,
    newBigintType,
    newDecimalType,
    newVarcharType,
    newTimestampType,
}

func NewType(id TypeID) Type {
    return typeConstructors[id]()
}

func GetTypeSize(id TypeID) (uint64, error) {
    switch (id) {
    case BOOLEAN:
    case TINYINT:
        return 1, nil
    case SMALLINT:
        return 2, nil
    case INTEGER:
        return 4, nil
    case BIGINT:
    case DECIMAL:
    case TIMESTAMP:
        return 8, nil
    case VARCHAR:
        return 0, nil
    default:
        break
    }
    return 0, errors.New("Unknown Type")
}

func TypeIDToString(id TypeID) string {
    switch (id) {
    case INVALID:
        return "INVALID"
    case BOOLEAN:
        return "BOOLEAN"
    case TINYINT:
        return "TINYINT"
    case SMALLINT:
        return "SMALLINT"
    case INTEGER:
      return "INTEGER"
    case BIGINT:
      return "BIGINT"
    case DECIMAL:
      return "DECIMAL"
    case TIMESTAMP:
      return "TIMESTAMP"
    case VARCHAR:
      return "VARCHAR"
    default:
      break;
    }
    return "INVALID"
}

func GetMinValue(id TypeID) *Value {
    switch (id) {
    case BOOLEAN:
        return NewValue(id, 0)
    case TINYINT:
        return NewValue(id, GOOSTUB_INT8_MIN)
    case SMALLINT:
        return NewValue(id, GOOSTUB_INT16_MIN)
    case INTEGER:
        return NewValue(id, GOOSTUB_INT32_MIN)
    case BIGINT:
        return NewValue(id, GOOSTUB_INT64_MIN)
    case DECIMAL:
        return NewValue(id, GOOSTUB_DECIMAL_MIN)
    case TIMESTAMP:
        return NewValue(id, 0)
    case VARCHAR:
        return NewValue(id, ([]byte)(nil))
    default:
        break
    }

    level.Error(common.Logger).Log("Can't get max value")
    return nil
}

func GetMaxValue(id TypeID) *Value {
    switch (id) {
    case BOOLEAN:
        return NewValue(id, 1)
    case TINYINT:
        return NewValue(id, GOOSTUB_INT8_MAX)
    case SMALLINT:
        return NewValue(id, GOOSTUB_INT16_MAX)
    case INTEGER:
        return NewValue(id, GOOSTUB_INT32_MAX)
    case BIGINT:
        return NewValue(id, GOOSTUB_INT64_MAX)
    case DECIMAL:
        return NewValue(id, GOOSTUB_DECIMAL_MAX)
    case TIMESTAMP:
        return NewValue(id, GOOSTUB_TIMESTAMP_MAX)
    case VARCHAR:
        return NewValue(id, ([]byte)(nil), false)
    default:
        break
    }

    level.Error(common.Logger).Log("Can't get max value")
    return nil
}

// singleton instances
var k_types = [14]Type {
    newInvalidType(),
    newBooleanType(),
    newTinyintType(),
    newSmallintType(),
    newIntegerType(),
    newBigintType(),
    newDecimalType(),
    newVarcharType(),
    newTimestampType(),
}

func GetInstance(id TypeID) Type {
    return k_types[id]
}
