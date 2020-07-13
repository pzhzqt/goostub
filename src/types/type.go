package types

import (
	"bytes"
	"common"
	"errors"
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

// The original comparison in BusTub is super WET, and thus I reduced it to a simple Compare and CompareTo function
// -1 => left < right, 0 => left == right, 1 => left > right
type CmpResult int

type Type interface {
	IsCoercableFrom(TypeID) bool
	GetTypeID() TypeID

	// Comparison
	Compare(*Value, *Value) (CmpResult, error)

	// Math Functions
	Add(*Value, *Value) (*Value, error)
	Subtract(*Value, *Value) (*Value, error)
	Multiply(*Value, *Value) (*Value, error)
	Divide(*Value, *Value) (*Value, error)
	Modulo(*Value, *Value) (*Value, error)
	Min(*Value, *Value) (*Value, error)
	Max(*Value, *Value) (*Value, error)
	Sqrt(*Value) (*Value, error)
	OperateNull(*Value, *Value) (*Value, error)

	IsZero(*Value) (bool, error)
	// Is the data in the struct storage or has indirection
	IsInlined(*Value) (bool, error)

	ToString(*Value) (string, error)
	SerializeTo(*Value, *bytes.Buffer) error
	DeserializeFrom(*bytes.Buffer) (*Value, error)
	Copy(*Value) *Value
	CastAs(*Value, TypeID) (*Value, error)
	// raw variable length data
	GetData(*Value) ([]byte, error)
	GetLength(*Value) (uint32, error)
}

func GetTypeSize(id TypeID) (uint64, error) {
	switch id {
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
	switch id {
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
		break
	}
	return "INVALID"
}

func GetMinValue(id TypeID) *Value {
	switch id {
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
	switch id {
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
var k_types = [14]Type{
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
