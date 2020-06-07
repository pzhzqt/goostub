package types

import (
    "fmt"
)

type TypeID int

const (
    INVALID TypeID = 0
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
    CmpFalse CmpBool = 0
    CmpTrue
    CmpNull
)

type Type interface {
    GetTypeSize(id TypeID)uint64
    IsCoercableFrom(id TypeID)bool
    GetTypeID()TypeID

    // Comparisons
    CompareEquals(Value, Value)CmpBool
    CompareNotEquals(Value, Value)CmpBool
    CompareLessThan(Value, Value)CmpBool
    CompareLessThanEquals(Value, Value)CmpBool
    CompareGreaterThan(Value, Value)CmpBool
    CompareGreaterThanEquals(Value, Value)CmpBool

    // Math Functions
    Add(Value, Value)Value
    Subtract(Value, Value)Value
    Multiply(Value, Value)Value
    Divide(Value, Value)Value
    Modulo(Value, Value)Value
    Min(Value, Value)Value
    Max(Value, Value)Value
    Sqrt(Value)Value
    OperateNull(Value, Value)Value
    IsZero(Value)Value
    // Is the data in the struct storage or has indirection
    IsInlined(Value)bool
    ToString(Value)string
    SerializeTo(Value, []byte)
    DeserializeFrom([]byte)Value
    Copy(Value)Value
    CastAs(Value, id TypeID)Value
    // raw variable length data
    GetData(Value)[]byte
    GetLength(Value)uint32
}

func TypeIDToString(id TypeID) {
}

func newInvalid()Type
func newBoolean()Type
func newTinyint()Type
func newSmallint()Type
func newInteger()Type
func newBigint()Type
func newDecimal()Type
func newVarchar()Type
func newTimestamp()Type

var typeConstructors = [...]func()Type {
    newInvalid,
    newBoolean,
    newTinyint,
    newSmallint,
    newInteger,
    newBigint,
    newDecimal,
    newVarchar,
    newTimestamp,
}

func NewType(id TypeID) Type {
    return typeConstructors[id]()
}

func GetMinValue(id TypeID) Value {
}

func GetMaxValue(id TypeID) Value {
}

// singleton instances
var k_types = [14]Type {
    newInvalid(),
    newBoolean(),
    newTinyint(),
    newSmallint(),
    newInteger(),
    newBigint(),
    newDecimal(),
    newVarchar(),
    newTimestamp(),
}

func GetInstance(id TypeID) Type {
    return k_types[id]
}
