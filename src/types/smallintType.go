package types

import (
	"log"
	"strconv"
	"unsafe"
)

type SmallintType struct {
	IntegerBaseType
}

func newSmallintType() *SmallintType {
	return &SmallintType{
		*newIntegerBaseType(TINYINT),
	}
}

func (t *SmallintType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "smallint_null", nil
	}

	val, ok := v.val.(int16)
	if !ok {
		log.Fatalln("smallint member function called on non-smallint value")
	}

	return strconv.FormatInt(int64(val), 10), nil
}

func (t *SmallintType) SerializeTo(v *Value, storage *byte) error {
	val, ok := v.val.(int16)
	if !ok {
		log.Fatalln("smallint member function called on non-smallint value")
	}

	*(*int16)(unsafe.Pointer(storage)) = val
	return nil
}

func (t *SmallintType) DeserializeFrom(storage *byte) (*Value, error) {
	val := *(*int16)(unsafe.Pointer(storage))
	return NewValue(t.id, val), nil
}
