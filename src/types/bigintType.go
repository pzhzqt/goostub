package types

import (
	"log"
	"strconv"
	"unsafe"
)

type BigintType struct {
	IntegerBaseType
}

func newBigintType() *BigintType {
	return &BigintType{
		*newIntegerBaseType(TINYINT),
	}
}

func (t *BigintType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "bigint_null", nil
	}

	val, ok := v.val.(int64)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	return strconv.FormatInt(int64(val), 10), nil
}

func (t *BigintType) SerializeTo(v *Value, storage *byte) error {
	val, ok := v.val.(int64)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	*(*int64)(unsafe.Pointer(storage)) = val
	return nil
}

func (t *BigintType) DeserializeFrom(storage *byte) (*Value, error) {
	val := *(*int64)(unsafe.Pointer(storage))
	return NewValue(t.id, val), nil
}
