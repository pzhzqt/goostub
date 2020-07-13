package types

import (
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
)

type IntegerType struct {
	IntegerBaseType
}

func newIntegerType() *IntegerType {
	return &IntegerType{
		*newIntegerBaseType(TINYINT),
	}
}

func (t *IntegerType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "integer_null", nil
	}

	val, ok := v.val.(int32)
	if !ok {
		log.Fatalln("integer member function called on non-integer value")
	}

	return strconv.FormatInt(int64(val), 10), nil
}

func (t *IntegerType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	val, ok := v.val.(int32)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	return binary.Write(storage, binary.LittleEndian, val)
}

func (t *IntegerType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	var val int32
	err := binary.Read(storage, binary.LittleEndian, &val)
	if err != nil {
		return nil, err
	}

	return NewValue(t.id, val), nil
}
