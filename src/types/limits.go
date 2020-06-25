package types

import (
	"math"
)

const (
	DBL_LOWEST = -math.MaxFloat64
	FLT_LOWEST = -math.MaxFloat32

	// Min values
	GOOSTUB_INT8_MIN      int8    = math.MinInt8 + 1
	GOOSTUB_INT16_MIN     int16   = math.MinInt16 + 1
	GOOSTUB_INT32_MIN     int32   = math.MinInt32 + 1
	GOOSTUB_INT64_MIN     int64   = math.MinInt64 + 1
	GOOSTUB_DECIMAL_MIN   float64 = FLT_LOWEST
	GOOSTUB_TIMESTAMP_MIN uint64  = 0
	GOOSTUB_DATE_MIN      uint32  = 0
	GOOSTUB_BOOLEAN_MIN   int8    = 0

	// Max values
	GOOSTUB_INT8_MAX      int8    = math.MaxInt8
	GOOSTUB_INT16_MAX     int16   = math.MaxInt16
	GOOSTUB_INT32_MAX     int32   = math.MaxInt32
	GOOSTUB_INT64_MAX     int64   = math.MaxInt64
	GOOSTUB_UINT64_MAX    uint64  = math.MaxUint64 - 1
	GOOSTUB_DECIMAL_MAX   float64 = math.MaxFloat64
	GOOSTUB_TIMESTAMP_MAX uint64  = 11231999986399999999
	GOOSTUB_DATE_MAX      uint64  = math.MaxInt32
	GOOSTUB_BOOLEAN_MAX   int8    = 1

	// Null values
	GOOSTUB_VALUE_NULL     uint32  = math.MaxUint32
	GOOSTUB_INT8_NULL      int8    = math.MinInt8
	GOOSTUB_INT16_NULL     int16   = math.MinInt16
	GOOSTUB_INT32_NULL     int32   = math.MinInt32
	GOOSTUB_INT64_NULL     int64   = math.MinInt64
	GOOSTUB_DATE_NULL      uint64  = 0
	GOOSTUB_TIMESTAMP_NULL uint64  = math.MaxUint64
	GOOSTUB_DECIMAL_NULL   float64 = DBL_LOWEST
	GOOSTUB_BOOLEAN_NULL   int8    = math.MinInt8

	// Max length
	GOOSTUB_VARCHAR_MAX_LEN uint32 = math.MaxUint32
	// TEXT = VARCHAR(GOOSTUB_TEXT_MAX_LEN)
	GOOSTUB_TEXT_MAX_LEN uint32 = 1000000000
	// VARCHAR objects WITH NULL length are NULL
	OBJECTLENGTH_NULL int = -1
)

func GetNull(id TypeID) interface{} {
	switch id {
	case INVALID:
		return GOOSTUB_VALUE_NULL
	case BOOLEAN:
		return GOOSTUB_BOOLEAN_NULL
	case TINYINT:
		return GOOSTUB_INT8_NULL
	case SMALLINT:
		return GOOSTUB_INT16_NULL
	case INTEGER:
		return GOOSTUB_INT32_NULL
	case BIGINT:
		return GOOSTUB_INT64_NULL
	case DECIMAL:
		return GOOSTUB_DECIMAL_NULL
	case TIMESTAMP:
		return GOOSTUB_TIMESTAMP_NULL
	default:
		break
	}

	return nil
}
