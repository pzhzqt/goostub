package types

import (
	"bytes"
	"common"
	"encoding/binary"
	"log"
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

func (t *VarcharType) GetLength(v *Value) (uint32, error) {
	return v.size, nil
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

	if l, _ := v.GetLength(); l == 0 {
		return "", nil
	}

	return string(v.val.([]byte)), nil
}

func (t *VarcharType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	val, ok := v.val.([]byte)
	if !ok {
		log.Fatalln("varchar member function called on non-varchar value")
	}
	err := binary.Write(storage, binary.LittleEndian, v.size)

	if v.size == GOOSTUB_VALUE_NULL || err != nil {
		return err
	}
	return binary.Write(storage, binary.LittleEndian, val)
	// TODO: left here; consider remove v.size since slice already has that
}
