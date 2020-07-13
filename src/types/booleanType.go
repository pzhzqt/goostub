package types

import (
	"bytes"
	"common"
	"encoding/binary"
	"log"
)

type BooleanType struct {
	BaseType
}

func newBooleanType() *BooleanType {
	return &BooleanType{
		*newBaseType(BOOLEAN),
	}
}

func (t *BooleanType) Compare(l *Value, r *Value) (CmpResult, error) {
	if l.GetTypeID() != BOOLEAN {
		log.Fatalln("BooleanType member function called from non-boolean type")
	}

	if !l.CheckComparable(r) || l.IsNull() || r.IsNull() {
		return 0, common.NewErrorf(common.MISMATCH_TYPE,
			"%s and %s are not comparable", TypeIDToString(l.typeID), TypeIDToString(r.typeID))
	}

	lval := l.val.(int8)
	rVal, err := r.CastAs(BOOLEAN)
	if err != nil {
		return 0, err
	}
	rval := rVal.val.(int8)

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

func (t *BooleanType) IsInlined(v *Value) (bool, error) {
	return true, nil
}

func (t *BooleanType) ToString(v *Value) (string, error) {
	if t.GetTypeID() != BOOLEAN {
		log.Fatalln("BooleanType member function called from non-boolean type")
	}

	if v.val.(int8) == 1 {
		return "true", nil
	} else if v.val.(int8) == 0 {
		return "false", nil
	}

	return "boolean_null", nil
}

func (t *BooleanType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	val, ok := v.val.(int8)
	if !ok {
		log.Fatalln("boolean member function called on non-boolean value")
	}

	return binary.Write(storage, binary.LittleEndian, val)
}

func (t *BooleanType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	var val int8
	err := binary.Read(storage, binary.LittleEndian, &val)
	if err != nil {
		return nil, err
	}
	return NewValue(BOOLEAN, val), nil
}

func (t *BooleanType) CastAs(v *Value, id TypeID) (*Value, error) {
	switch id {
	case BOOLEAN:
		return t.Copy(v), nil
	case VARCHAR:
		if v.IsNull() {
			return NewValue(VARCHAR, nil, false), nil
		}
		s, _ := t.ToString(v)
		return NewValue(VARCHAR, s), nil
	default:
		break
	}
	return nil, common.NewErrorf(common.INVALID,
		"Boolean is not coearcible to %s", TypeIDToString(id))
}
