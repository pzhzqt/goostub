package types

import (
	"bytes"
	"common"
	"encoding/binary"
	"log"
	"strconv"
	"strings"
)

type VarcharType struct {
	BaseType
}

func newVarcharType() *VarcharType {
	return &VarcharType{
		*newBaseType(VARCHAR),
	}
}

func (t *VarcharType) GetData(v *Value) ([]byte, error) {
	return v.val.([]byte), nil
}

func (t *VarcharType) GetLength(v *Value) int32 {
	val, ok := v.val.([]byte)
	if !ok {
		return -1
	}

	return int32(len(val))
}

func (t *VarcharType) Compare(l, r *Value) (CmpResult, error) {
	if l.GetTypeID() != VARCHAR {
		log.Fatalln("Varchar function called by non-varchar value")
	}

	if l.IsNull() || r.IsNull() {
		return 0, common.NewError(common.INVALID, "Null Value is not comparable")
	}

	lval := l.val.([]byte)
	if r.GetTypeID() != VARCHAR {
		r, _ = r.CastAs(VARCHAR)
	}
	rval := r.val.([]byte)

	return CmpResult(bytes.Compare(lval, rval)), nil
}

func (t *VarcharType) Min(l, r *Value) (*Value, error) {
	res, err := t.Compare(l, r)
	if err != nil {
		return nil, err
	}

	if res < 0 {
		return l.Copy(), nil
	}

	return r.Copy(), nil
}

func (t *VarcharType) Max(l, r *Value) (*Value, error) {
	res, err := t.Compare(l, r)
	if err != nil {
		return nil, err
	}

	if res > 0 {
		return l.Copy(), nil
	}

	return r.Copy(), nil
}

func (t *VarcharType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "varchar_null", nil
	}

	if v.GetLength() == 0 {
		return "", nil
	}

	return string(v.val.([]byte)), nil
}

func (t *VarcharType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	if v.isNull {
		return binary.Write(storage, binary.LittleEndian, GOOSTUB_VALUE_NULL)
	}
	val, ok := v.val.([]byte)
	if !ok {
		log.Fatalln("varchar member function called on non-varchar value")
	}
	err := binary.Write(storage, binary.LittleEndian, int32(len(val))) // int is not fixed-length type, need int32
	if err != nil {
		return err
	}
	return binary.Write(storage, binary.LittleEndian, val)
}

func (t *VarcharType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	var l int32
	if err := binary.Read(storage, binary.LittleEndian, &l); err != nil {
		return nil, err
	}
	data := make([]byte, l)
	if err := binary.Read(storage, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	// set manage_data as true
	return NewValue(VARCHAR, data), nil
}

func (t *VarcharType) CastAs(v *Value, id TypeID) (*Value, error) {
	var str string
	var maxVal, minVal int64
	var bitSize int

	switch id {
	case BOOLEAN:
		str = strings.ToLower(v.String())
		if str == "true" || str == "1" || str == "t" {
			return NewValue(id, 1), nil
		}
		if str == "false" || str == "0" || str == "f" {
			return NewValue(id, 0), nil
		}
		return nil, common.NewError(common.CONVERSION, "Boolean value format error")
	case TINYINT:
		maxVal = int64(GOOSTUB_INT8_MAX)
		minVal = int64(GOOSTUB_INT8_MIN)
		bitSize = 8
		goto MakeInt
	case SMALLINT:
		maxVal = int64(GOOSTUB_INT16_MAX)
		minVal = int64(GOOSTUB_INT16_MIN)
		bitSize = 16
		goto MakeInt
	case INTEGER:
		maxVal = int64(GOOSTUB_INT32_MAX)
		minVal = int64(GOOSTUB_INT32_MIN)
		bitSize = 32
		goto MakeInt
	case BIGINT:
		maxVal = int64(GOOSTUB_INT64_MAX)
		minVal = int64(GOOSTUB_INT64_MIN)
		bitSize = 64
		goto MakeInt
	case DECIMAL:
		str = v.String()
		flval, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		if flval > GOOSTUB_DECIMAL_MAX || flval < GOOSTUB_DECIMAL_MIN {
			return nil, common.NewError(common.OUT_OF_RANGE, "Numeric value out of range")
		}
		return NewValue(id, flval), nil
	case TIMESTAMP:
		str = v.String()
		intval, err := strconv.ParseUint(str, 0, 64)
		if err != nil {
			return nil, err
		}
		if intval > GOOSTUB_TIMESTAMP_MAX || intval < GOOSTUB_TIMESTAMP_MIN {
			return nil, common.NewError(common.OUT_OF_RANGE, "Timestamp value out of range")
		}
		return NewValue(id, intval), nil
	case VARCHAR:
		return v.Copy(), nil
	default:
		break
	}

	return nil, common.NewErrorf(common.MISMATCH_TYPE, "Varchar is not coercable to %s", TypeIDToString(id))

MakeInt:
	str = v.String()
	intval, err := strconv.ParseInt(str, 0, bitSize)
	if err != nil {
		return nil, err
	}
	if intval > maxVal || intval < minVal {
		return nil, common.NewError(common.OUT_OF_RANGE, "Numeric value out of range")
	}
	return NewValue(id, intval), nil
}
