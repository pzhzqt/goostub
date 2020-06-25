package types

import (
	"log"
	"strconv"
	"unsafe"
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

func (t *IntegerType) SerializeTo(v *Value, storage *byte) error {
	val, ok := v.val.(int32)
	if !ok {
		log.Fatalln("integer member function called on non-integer value")
	}

	*(*int32)(unsafe.Pointer(storage)) = val
	return nil
}

func (t *IntegerType) DeserializeFrom(storage *byte) (*Value, error) {
	val := *(*int32)(unsafe.Pointer(storage))
	return NewValue(t.id, val), nil
}
