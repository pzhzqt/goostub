package types

import (
	"bytes"
	"encoding/binary"
	"goostub/common"
	"log"
	"time"
)

type TimestampType struct {
	BaseType
}

func newTimestampType() *TimestampType {
	return &TimestampType{
		*newBaseType(TIMESTAMP),
	}
}

func (t *TimestampType) Min(l, r *Value) (*Value, error) {
	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	res, err := t.Compare(l, r)
	if err != nil {
		return nil, err
	}
	if res < 0 {
		return l.Copy(), nil
	}

	return r.Copy(), nil
}

func (t *TimestampType) Max(l, r *Value) (*Value, error) {
	if l.IsNull() || r.IsNull() {
		return l.OperateNull(r)
	}

	res, err := t.Compare(l, r)
	if err != nil {
		return nil, err
	}
	if res > 0 {
		return l.Copy(), nil
	}

	return r.Copy(), nil
}

func (t *TimestampType) Compare(l, r *Value) (CmpResult, error) {
	if l.GetTypeID() != TIMESTAMP {
		log.Fatalln("Time stamp member function called on non-timestamp value")
	}

	if !l.CheckComparable(r) {
		return 0, common.NewErrorf(common.MISMATCH_TYPE,
			"Mismatched types of %s and %s", TypeIDToString(l.typeID), TypeIDToString(r.typeID))
	}

	if l.IsNull() || r.IsNull() {
		return 0, common.NewError(common.INVALID,
			"Null Value is not comparable")
	}

	lval := r.val.(uint64)
	rval := r.val.(uint64)

	var ret CmpResult

	if lval < rval {
		ret = -1
	} else if lval > rval {
		ret = 1
	} else {
		ret = 0
	}

	return ret, nil

}

func (t *TimestampType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "timestamp_null", nil
	}
	return time.Unix(0, int64(v.val.(uint64))).String(), nil
}

func (t *TimestampType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	val, ok := v.val.(uint64)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	return binary.Write(storage, binary.LittleEndian, val)
}

func (t *TimestampType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	var val uint64
	err := binary.Read(storage, binary.LittleEndian, &val)
	if err != nil {
		return nil, err
	}

	return NewValue(t.id, val), nil
}

func (t *TimestampType) CastAs(v *Value, id TypeID) (*Value, error) {
	switch id {
	case TIMESTAMP:
		return t.Copy(v), nil
	case VARCHAR:
		if v.IsNull() {
			return NewValue(VARCHAR, nil, false), nil
		}
		s, _ := v.ToString()
		return NewValue(VARCHAR, s), nil
	default:
		break
	}

	return nil, common.NewErrorf(common.CONVERSION, "Timestamp is not coercable to %s", TypeIDToString(v.GetTypeID()))
}
